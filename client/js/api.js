let config = {
    mode: "string",
    length: 6,
    noise: 0,
    lines: 3,
    width: 240,
    height: 60,
    fonts: [
        "vogue.ttf",
    ],
    bg: {
        r: 255,
        g: 255,
        b: 255,
        a: 0
    },
    source: "ABCDEFGHJKMNOQRSTUVXYZabcdefghjkmnoqrstuvxyz123456789",
    skew: 0,
    dots: 0,
}

let currentHash = 0

getCustom()

function getCustom() {
    var xmlHttp = new XMLHttpRequest()
    xmlHttp.open("POST", "/custom/", false)
    xmlHttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    xmlHttp.send(JSON.stringify(config))

    const obj = JSON.parse(xmlHttp.responseText)
    const url = obj.captcha
    currentHash = obj.hash

    document.querySelector("#captcha-path a").textContent = obj.hash
    document.querySelector("#captcha-path a").href = url

    xmlHttp.open("GET", url, false)
    xmlHttp.send(null)

    document.getElementById("captcha-image").innerHTML = xmlHttp.responseText
}

function generate() {
    getCustom()
}

function verify() {
    var xmlHttp = new XMLHttpRequest()
    xmlHttp.open("POST", "/verify/", false)
    xmlHttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");

    let data = {
        hash: currentHash,
        code: document.querySelector('.captcha-form #input-form').value
    }

    xmlHttp.send(JSON.stringify(data))

    const obj = JSON.parse(xmlHttp.responseText)

    if (obj.verify == "true") {
        alert("Verification has passed")
    } else {
        alert("Verification failed")
    }
}

function changeMenuToString() {

    document.querySelector('.check-boxes').style.display="flex"
    document.querySelector('.input-forms').style.display="flex"

    document.querySelector('#noise').style.display = "block"
    document.querySelector('#lines').style.display = "block"

    document.querySelector('#skew').style.display="none"
    document.querySelector('#dots').style.display="none"

    document.querySelector('.input-forms #source').style.display="block"

    document.querySelector('#length').style.display = "block"

    config.mode="string"
    getCustom()
}

function changeMenuToDigits() {
    // Height   int
	// Width    int
	// Length   int
	// MaxSkew  float64
	// DotCount int

    // check-boxes

    // input-forms

    document.querySelector('.check-boxes').style.display="none"
    document.querySelector('.input-forms').style.display="none"

    document.querySelector('#noise').style.display = "none"
    document.querySelector('#lines').style.display = "none"

    document.querySelector('#skew').style.display="block"
    document.querySelector('#dots').style.display="block"

    config.mode="digits"
    getCustom()
}

function changeMenuToMath() {

    document.querySelector('.check-boxes').style.display="flex"
    document.querySelector('.input-forms').style.display="flex"

    document.querySelector('#noise').style.display = "block"
    document.querySelector('#lines').style.display = "block"

    document.querySelector('#skew').style.display="none"
    document.querySelector('#dots').style.display="none"

    document.querySelector('#length').style.display = "none"
    document.querySelector('.input-forms #source').style.display = "none"

    config.mode="math"
    getCustom()

    // range-noise-wrap
    // .input-forms #source
    // .ragnes #range-length-wrap
}

function changeMenu() {

    var r1 = document.querySelector('.menu #r1').checked
    var r2 = document.querySelector('.menu #r2').checked
    var r3 = document.querySelector('.menu #r3').checked

    if (r1) {
        changeMenuToMath()
    } else if (r2) {
        changeMenuToString()
    } else if (r3) {
        changeMenuToDigits()
    }
}

function addFont(fontName, element) {
    if (config.fonts.includes(fontName)) {
        config.fonts.splice(config.fonts.indexOf(fontName), 1)
    } else {
        config.fonts.push(fontName)
    }
    getCustom()
}

function rangeLength(newVal) {
    config.length = +newVal
    document.querySelector('#range-length').textContent = newVal
    getCustom()
}

function rangeNoise(newVal) {
    config.noise = +newVal
    document.querySelector('#range-noise').textContent = newVal
    getCustom()
}

function rangeLines(newVal) {
    config.lines = +newVal
    document.querySelector('#range-lines').textContent= newVal
    getCustom()
}

function rangeWidth(newVal) {
    config.width = +newVal
    document.querySelector('#range-width').textContent = newVal
    getCustom()
}

function rangeHeight(newVal) {
    config.height = +newVal
    document.querySelector('#range-height').textContent = newVal
    getCustom()
}

function rangeSkew(newVal) {
    config.skew = +newVal
    document.querySelector('#range-skew').textContent = newVal
    getCustom()
}

function rangeDots(newVal) {
    config.dots = +newVal
    document.querySelector('#range-dots').textContent = newVal
    getCustom()
}