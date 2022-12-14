package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/glycosupport/gocaptcha-service/pkg/store"
	"github.com/glycosupport/gocaptcha-service/pkg/utils"
)

var addr string

type captchaServer struct {
	store *store.CaptchaStore
}

func NewCaptchaServer() *captchaServer {
	return &captchaServer{store: store.New()}
}

func (cs *captchaServer) generateCaptchaHandler(c *gin.Context) {

	data, err := cs.store.GenerateCaptcha(addr)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": data.Code, "captcha": data.URL})
}

func (cs *captchaServer) generateCustomCaptchaHandler(c *gin.Context) {

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	defer c.Request.Body.Close()

	var request store.CaptchaRequest

	json.Unmarshal(jsonData, &request)

	data, err := cs.store.GenerateCustomCaptcha(addr, &request)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": data.Code, "captcha": data.URL, "hash": data.Hash})
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
	} else {
		c.JSON(http.StatusOK, gin.H{"verify": "false"})
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

func (cs *captchaServer) removeCaptchaHandler(c *gin.Context) {

	name := c.Param("name")

	e := os.Remove("./captchas/" + name)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"remove": "false"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"remove": "true"})
}

func main() {

	if _, err := os.Stat("./captchas/"); os.IsNotExist(err) {
		if err := os.Mkdir("./captchas/", 0777); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}

	if _, err := os.Stat("./logs/"); os.IsNotExist(err) {
		if err := os.Mkdir("./logs/", 0777); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}

	file, err := os.OpenFile("logs/common.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	router := gin.Default()
	server := NewCaptchaServer()

	ip := os.Getenv("IP")
	port := os.Getenv("PORT")

	if ip == "" {
		log.Println("IP address not set, local is used")
		ip = utils.GetLocalIP()
	}

	if port == "" {
		log.Println("PORT not set, 8080 is used")
		port = "8080"
	}

	addr = ip + ":" + port

	gin.SetMode(gin.DebugMode)

	router.Use(static.Serve("/", static.LocalFile("./client/", true)))
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	router.POST("/custom", server.generateCustomCaptchaHandler)
	router.POST("/verify", server.verifyCaptchaHandler)
	router.GET("/remove/:name", server.removeCaptchaHandler)

	router.GET("/generate", server.generateCaptchaHandler)
	router.GET("/:name", server.getCaptchaHandler)

	router.Run(addr)
}
