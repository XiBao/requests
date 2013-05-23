# Requests

[https://github.com/XiBao/requests](https://github.com/XiBao/requests)

Go HTTP Requests for Humans™, inspired by [github.com/kennethreitz/requests](https://github.com/kennethreitz/requests)

## Usage

```
package main

import (
	"github.com/XiBao/requests"
	"log"
	"net/url"
	"time"
	"net/http"
)

func testLogin() {
	session := requests.NewSession()
	res, err := session.Get(&requests.Request{Url:"http://sem.xibao100.com/user/signin"})
	log.Println(res.Header, res.EffectiveUrl, res.Redirects)
	params := &url.Values{}
	params.Add("uname", "syd")
	params.Add("passwd", "xxxxxxx")
	request := &requests.Request{
		Url: "http://sem.xibao100.com/user/signin",
		Params: params,
	}
	res, err = session.Post(request)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Header, res.EffectiveUrl, res.Redirects)
}

func testEtagAndLastModified() {
	//reqUrl := "http://dnn506yrbagrg.cloudfront.net/pages/scripts/0013/7219.js?380363"
	reqUrl := "https://ssl.google-analytics.com/ga.js"
	session := requests.NewSession()
	res, err := session.Get(&requests.Request{Url:reqUrl})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Header, res.EffectiveUrl)
	lastModified, err := time.Parse(http.TimeFormat, res.Header.Get("Last-Modified"))
	Etag := res.Header.Get("Etag")
	request := &requests.Request{
		Url: reqUrl,
		Etag: Etag,
		LastModified: lastModified,
	}
	res, err = session.Get(request)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Header, res.EffectiveUrl, res.StatusCode)
}

func main() {
	testEtagAndLastModified()
	//testLogin()
}
```