package requests

import (
	"net/url"
	"net/http"
	"strings"
	"io/ioutil"
	"github.com/salviati/go-guess/guess"
	"time"
	"github.com/bububa/cookiejar"
	"sync"
	"regexp"
	"errors"
	"fmt"
	"math/rand"
	"bytes"
	"mime/multipart"
	"os"
	"io"
	//"log"
)

func Get(request *Request) (response *Response, err error) {
	request.Method = GET_METHOD
	session := NewSession()
	return session.Do(request)
}

func Post(request *Request) (response *Response, err error) {
	request.Method = POST_METHOD
	session := NewSession()
	return session.Do(request)
}

func Put(request *Request) (response *Response, err error) {
	request.Method = PUT_METHOD
	session := NewSession()
	return session.Do(request)
}

func Head(request *Request) (response *Response, err error) {
	request.Method = HEAD_METHOD
	session := NewSession()
	return session.Do(request)
}

func NewSession() *Session {
	session := &Session{
		Client:&http.Client{Jar: cookiejar.NewJar(true)},
		Mutex:new(sync.Mutex),
		Redirects:make(map[string][]string),
		MaxRedirect:10,
	}
	session.Client.CheckRedirect = session.checkRedirect
	return session
}

func (this *Session) SetMaxRedirect(maxRedirect int) {
	this.MaxRedirect = maxRedirect
}

func (this *Session) Get(request *Request) (response *Response, err error) {
	request.Method = GET_METHOD
	return this.Do(request)
}

func (this *Session) Post(request *Request) (response *Response, err error) {
	request.Method = POST_METHOD
	return this.Do(request)
}

func (this *Session) Put(request *Request) (response *Response, err error) {
	request.Method = PUT_METHOD
	return this.Do(request)
}

func (this *Session) Head(request *Request) (response *Response, err error) {
	request.Method = HEAD_METHOD
	return this.Do(request)
}

func (this *Session) Do(request *Request) (response *Response, err error) {
	this.Mutex.Lock()
	defer func() {
		this.Mutex.Unlock()
	}()

	var httpRequest *http.Request
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	if request.Params == nil {
		request.Params = &url.Values{}
		httpRequest, err = http.NewRequest(request.Method, request.Url, nil)
	}else if request.Method == POST_METHOD {
		if len(request.Files) > 0 {
			body_buf := bytes.NewBufferString("")
	        body_writer := multipart.NewWriter(body_buf)
	        for k, vs := range *request.Params {
	            for _, v := range vs {
	            	body_writer.WriteField(k, v)
	            }
	        }
	        for _, formFile := range request.Files {
	            file_writer, err := body_writer.CreateFormFile(formFile.Name, formFile.Filename)
	            if err != nil {
	            	continue
	            }
	            if formFile.Data == nil {
	            	fh, err := os.Open(formFile.Filename)
		            if err != nil {
		            	continue
		            }
		            io.Copy(file_writer, fh)
	            }else{
	            	io.Copy(file_writer, bytes.NewReader(formFile.Data))
	            }
	            
	        }
	        request.Headers["Content-Type"] = body_writer.FormDataContentType()
	        body_writer.Close()
	        httpRequest, err = http.NewRequest(request.Method, request.Url, body_buf)
		}else{
			request.Headers["Content-Type"] = "application/x-www-form-urlencoded"
			httpRequest, err = http.NewRequest(request.Method, request.Url, strings.NewReader(request.Params.Encode()))
		}
	}
	if err != nil {
		return nil, err
	}
	if request.BasicAuth != nil {
		httpRequest.SetBasicAuth(request.BasicAuth.Username, request.BasicAuth.Password)
	}
	if request.UserAgent != "" {
		request.Headers["UserAgent"] = request.UserAgent
	}

	if request.Etag != "" {
		request.Headers["If-None-Match"] = request.Etag
	}
	if !request.LastModified.IsZero() {
		request.Headers["If-Modified-Since"] = request.LastModified.Format(http.TimeFormat)
	}
	if len(request.Headers) > 0 {
		for k, v := range request.Headers {
			httpRequest.Header.Set(k, v)
		}
	}
	request.Url = httpRequest.URL.String()
	if this.Client.Jar != nil {
        for _, cookie := range this.Client.Jar.Cookies(httpRequest.URL) {
            httpRequest.AddCookie(cookie)
        }
    }
    startAt := time.Now()
    transport := &http.Transport { ResponseHeaderTimeout: request.ConnectionTimeout }
    if len(this.Proxies) > 0 {
    	proxyUrl, err := GetRandomProxy(this.Proxies)
    	if err == nil {
    		transport.Proxy = http.ProxyURL(proxyUrl)
    	}
    }
    this.Client.Transport = transport
    httpResponse, err := this.Client.Do(httpRequest)
    if err != nil {
        return nil, err
    }
    defer httpResponse.Body.Close()
    body, err := ioutil.ReadAll(httpResponse.Body)
    if err != nil {
        return nil, err
    }

    var encoding string
    contentType := httpResponse.Header.Get("Content-Type")
    if len(contentType) > 0 {
    	pattern := regexp.MustCompile(`charset\=([^;,\r\n]*)`)
    	matched := pattern.FindStringSubmatch(contentType)
    	if len(matched) > 1 {
    		encoding = matched[1]
    	}
    }
    if encoding == "" && len(body) > 0{
    	guess.DetermineEncoding(body, guess.CN)
    }else{
    	encoding = "utf-8"
    }

    response = &Response{
    	Url: request.Url,
    	Encoding: encoding,
    	Content: body,
    	Cookies: this.Client.Jar.Cookies(httpRequest.URL),
    	StartAt: startAt,
    	Elapsed: time.Now().Sub(startAt),
    	StatusCode: httpResponse.StatusCode,
    	Header: httpResponse.Header,
    }

    if eurl := httpResponse.Header.Get("Url"); eurl != "" {
    	response.EffectiveUrl = eurl
    }else if redirects, found := this.Redirects[request.Url]; found && len(redirects)>0 {
    	response.EffectiveUrl = redirects[0]
    	response.Redirects = redirects
    	delete(this.Redirects, request.Url)
    }else {
    	response.EffectiveUrl = request.Url
    }

    return response, nil
}

func (this *Session) checkRedirect(req *http.Request, via []*http.Request) error {
	originalUrl := via[0].URL.String()
	this.Redirects[originalUrl] = append(this.Redirects[originalUrl], req.URL.String())
	if len(via) >= this.MaxRedirect {
		return errors.New(fmt.Sprintf("stopped after %d redirects", this.MaxRedirect))
	}
	return nil
}

func GetRandomProxy(arr []string) (proxyUrl *url.URL, err error) {
    idx := rand.Intn(len(arr))
    proxy := arr[idx]
    if !strings.HasPrefix(proxy, "http") { proxy = "http://" + proxy }
    proxyUrl, err = url.Parse(proxy)
    return
}