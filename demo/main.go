package main

import (
	"github.com/XiBao/requests"
	"log"
	"net/url"
)

func main() {
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