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

function changeMenu() {

    var r1 = document.querySelector('.menu #r1')
    var r2 = document.querySelector('.menu #r2')
    var r3 = document.querySelector('.menu #r3')

    console.log(r1.checked, r2.checked, r3.checked)

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
    getCustom()
}

function rangeNoise(newVal) {
    config.noise = +newVal
    getCustom()
}

function rangeLines(newVal) {
    config.lines = +newVal
    getCustom()
}

function rangeWidth(newVal) {
    config.width = +newVal
    getCustom()
}

function rangeHeight(newVal) {
    config.height = +newVal
    getCustom()
}

function rangeSkew(newVal) {
    config.skew = +newVal
    getCustom()
}

function rangeDots(newVal) {
    config.dots = +newVal
    getCustom()
}