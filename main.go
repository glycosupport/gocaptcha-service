package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/glycosupport/gocaptcha-service/docs"
	"github.com/glycosupport/gocaptcha-service/pkg/store"
	"github.com/glycosupport/gocaptcha-service/pkg/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var addr string

type captchaServer struct {
	store *store.CaptchaStore
}

func NewCaptchaServer() *captchaServer {
	return &captchaServer{store: store.New()}
}

func (cs *captchaServer) generateCaptchaHandler(c *gin.Context) {

	ip := utils.GetLocalIP()

	if len(ip) == 0 {
		ip = "localhost"
	}

	ip += ":8080" // port

	data, err := cs.store.GenerateCaptcha(ip)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": data.Code, "captcha": data.URL})
}

func (cs *captchaServer) generateCustomCaptchaHandler(c *gin.Context) {

	ip := utils.GetLocalIP()

	if len(ip) == 0 {
		ip = "localhost"
	}

	ip += ":8080" // port

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	defer c.Request.Body.Close()

	var request store.CaptchaRequest

	json.Unmarshal(jsonData, &request)

	fmt.Println(request)

	data, err := cs.store.GenerateCustomCaptcha(ip, &request)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": data.Code, "captcha": data.URL})
}

func (cs *captchaServer) verifyCaptchaHandler(c *gin.Context) {

	type RequestVerify struct {
		Hash string `json:hash`
		Code string `json:code`
	}

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	defer c.Request.Body.Close()

	var request RequestVerify

	json.Unmarshal(jsonData, &request)

	if cs.store.VerifyCaptcha(request.Hash, request.Code) {
		c.JSON(http.StatusOK, gin.H{"verify": "true"})
		fmt.Println("VERIFY TRUE")
	} else {
		c.JSON(http.StatusOK, gin.H{"verify": "false"})
		fmt.Println("VERIFY FALSE")
	}
}

func (cs *captchaServer) getCaptchaHandler(c *gin.Context) {

	b64s, err := cs.store.GetCaptcha(c.Param("name"))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	returnValue := "<img style=\"display: block;\" src=\""
	returnValue += b64s + "\">"

	c.Data(200, "text/html; charset=utf-8", []byte(returnValue))
}

func (cs *captchaServer) getStaticHandler(c *gin.Context) {
	file, err := os.Open("./static/index.html")

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	c.Data(200, "text/html; charset=utf-8", []byte(bytes))
}

func main() {

	docs.SwaggerInfo.Title = "GO-CAPTHCA SERVICE API"
	docs.SwaggerInfo.Description = "Captcha generation service."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router := gin.Default()
	server := NewCaptchaServer()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", server.getStaticHandler)
	router.POST("/custom", server.generateCustomCaptchaHandler)
	router.POST("/verify", server.verifyCaptchaHandler)
	router.GET("/generate", server.generateCaptchaHandler)
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.GET("/:name", server.getCaptchaHandler)

	router.Run()
}
