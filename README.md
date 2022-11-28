
# gocaptcha-service :green_heart:

Captcha generation and verification service in Go

- [Installation](https://github.com/glycosupport/gocaptcha-service#installation)

- [Run Locally](https://github.com/glycosupport/gocaptcha-service#run-locally)

- [Docker Images](https://github.com/glycosupport/gocaptcha-service#docker-images)

- [API Reference](https://github.com/glycosupport/gocaptcha-service#api-reference)

- [Description](https://github.com/glycosupport/gocaptcha-service#description)

- [Screenshots](https://github.com/glycosupport/gocaptcha-service#screenshots)

- [Tech Stack](https://github.com/glycosupport/gocaptcha-service#tech-stack)

## Installation

Install gocaptcha-service

```bash
  git clone https://github.com/glycosupport/gocaptcha-service
  cd gocaptcha-service
```
    
## Run Locally

Add environment variables `if necessary`

If not added, the default settings will be used (local address and 8080 port)

```bash
   export PORT=8080
   export IP=127.0.0.1
```

Get dependencies

```bash
   go mod download
```

Start the server

```bash
   go run main.go
```

Open with a browser

```
   http://IP:PORT/
```

## Docker Images

#### Build development image:

```bash
   sudo docker build -t gocaptcha-service .
```

#### Run container:

```bash
   sudo docker run -it --rm gocaptcha-service
```

## API Reference

#### Generate default captcha

```
  GET /generate/
```

Response:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `captcha` | `string` | URL where can get the captcha image  |
| `code`    | `string` | Code captcha |

#### Generate custom captcha

```
  POST /custom/
```

Request:

| Parameter     | Type       | Description                       |
| :--------     | :-------   | :-------------------------------- |
| `mode`        | `string`   | "string", "digits" or "math" |
| `length`      | `int`      | Captcha word length |
| `noise`       | `int`      | Noise level |
| `lines`       | `int`      | Line noise level |
| `width`       | `int`      | Width captcha image |
| `height`      | `int`      | Height captcha image |
| `fonts`       | `string[]` | Array captcha fonts |
| `bg`          | `object`   | Colors RGBA background |
| `source`      | `string`   | Captcha alphabet |
| `skew`        | `float`    | Max skew [digits]|
| `dots`        | `int`      | Max count dots [digits]|


Response:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `captcha` | `string` | URL where can get the captcha image |
| `code`    | `string` | Code captcha |

#### Get captcha image

```
  GET /:name
```

Example: http://localhost:8080/5a5e1f5ecc6d0b8ac4443172561d8acb.png

#### Verify captcha

```
  POST /verify/
```

Request:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `hash` | `string` | Hash of necessary captcha |
| `code`    | `string` | Assumed code |

Response:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `verify` | `string` | "false" or "true" |

#### Remove captcha from server

```
  GET /remove/:name
```

Example: http://localhost:8080/remove/5a5e1f5ecc6d0b8ac4443172561d8acb.png


## Description

Text captchas are the most commonly used, but because of their simple structure, they are also the most vulnerable. This type of captcha, as a rule, requires recognizing a geometrically distorted sequence consisting of letters or numbers.

To increase security, various protection mechanisms are used, which can be divided into anti-segmentation and anti-identification. The first group of mechanisms is aimed at worsening the process of highlighting individual characters, the second – at recognizing the characters themselves.

Methods of protection:

[Anti-Segmentation]

- Hollow symbols

- CCT and overlap

- Background noise

- Two-level structure


[Anti-identification]

- Random fonts

- Random turns

- Wavelike symbols

- Various languages

This service uses bundles of various methods.

The following security methods are used in the "string" module:

- Background noises [generating different lines and backgrounds from random characters]

- Random fonts

- Random text length


The following security methods are used in the "digits" module:

- CCT and overlap

- Background noise

- Hollow symbols

- Wavelike symbols

- Random turns


The "math" module uses the same security methods as the "string" module. For additional protection, a method is used where it is necessary to keep the result of a mathematical expression instead of the captcha code.

Main structure captcha package:

```go
type Captcha struct {
    Driver Driver
}
```

A separate driver is used for each captcha module with the main functions *DrawCaptcha* and *GenerateIdQuestionAnswer*

- **Module digits**

    - First, a random position for the captcha is generated
    - Each digit is being drawn
    - Noise is applied using distortion lines
    - Filling with noise in the form of circles according to a given parameter

```go
func (d *DriverDigit) DrawCaptcha(content string) (item Item, err error) {
    itemDigit := NewItemDigit(d.Width, d.Height, d.DotCount, d.MaxSkew)
    digits := stringToFakeByte(content)

    itemDigit.calculateSizes(d.Width, d.Height, len(digits))
    maxx := d.Width - (itemDigit.width+itemDigit.dotSize)*len(digits) - itemDigit.dotSize
    maxy := d.Height - itemDigit.height - itemDigit.dotSize*2

    var border int

    if d.Width > d.Height {
        border = d.Height / 5
    } else {
        border = d.Width / 5
    }

    x := rand.Intn(maxx-border*2) + border
    y := rand.Intn(maxy-border*2) + border

    for _, n := range digits {
        itemDigit.drawDigit(digitFontData[n], x, y)
        x += itemDigit.width + itemDigit.dotSize
    }

    itemDigit.strikeThrough()
    itemDigit.distort(rand.Float64()*(10-5)+5, rand.Float64()*(200-100)+100)
    itemDigit.fillWithCircles(d.DotCount, itemDigit.dotSize)

    return itemDigit, nil
}
```
- **Module string**
    - First, a random position for the captcha is generated
    - Each digit is being drawn
    - Noise is applied using distortion lines
    - Filling with noise in the form of circles according to a given parameter

```go
func (d *DriverString) DrawCaptcha(content string) (item Item, err error) {

    var bgc color.RGBA

    if d.BgColor != nil {
        bgc = *d.BgColor
    } else {
        bgc = RandLightColor()
    }

    itemChar := NewItemChar(d.Width, d.Height, bgc)

    if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
        itemChar.drawHollowLine()
    }

    if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
        itemChar.drawSlimLine(3)
    }

    if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
        itemChar.drawSineLine()
    }

    if d.NoiseCount > 0 {
        source := TxtNumbers + TxtAlphabet + ",.[]<>"
        noise := RandText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
        err = itemChar.drawNoise(noise, d.fontsArray)
        if err != nil {
            return
        }
    }

    err = itemChar.drawText(content, d.fontsArray)
    if err != nil {
        return
    }

    return itemChar, nil
}
```
- **Module math**

    - Same methods are used as in the "string" module
    - Logic of calculating a mathematical expression

```go
func (d *DriverMath) DrawCaptcha(question string) (item Item, err error) {
    var bgc color.RGBA
    
    if d.BgColor != nil {
        bgc = *d.BgColor
    } else {
        bgc = RandLightColor()
    }

    itemChar := NewItemChar(d.Width, d.Height, bgc)

    if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
        itemChar.drawHollowLine()
    }

    if d.NoiseCount > 0 {
        noise := RandText(d.NoiseCount, strings.Repeat(TxtNumbers, d.NoiseCount))
        err = itemChar.drawNoise(noise, fontsAll)
        if err != nil {
            return
        }
    }

    if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
        itemChar.drawSlimLine(3)
    }

    if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
        itemChar.drawSineLine()
    }

    err = itemChar.drawText(question, d.fontsArray)
    if err != nil {
        return
    }
    return itemChar, nil
}
```

## Screenshots

#### **ROOT URL**

![Main Frame](https://raw.githubusercontent.com/glycosupport/gocaptcha-service/dev/screenshots/frame.png)

### **GIN Requests**

![Gin Requests](https://raw.githubusercontent.com/glycosupport/gocaptcha-service/dev/screenshots/gin.png)

## Tech Stack

**Client:** HTML, CSS, JS

**Server:** Golang, Gin, Swag


