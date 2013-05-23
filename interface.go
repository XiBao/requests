package requests

import (
	"net/url"
    "net/http"
    "time"
    "sync"
)

type BasicAuth struct {
	Username string
	Password string
}

type FormFile struct {
	Name string
	Filename string
	Data []byte
}

type Request struct {
	Method string
	Url string
	Headers map[string]string
	Files []FormFile
	Params *url.Values
	RawData string
	BasicAuth *BasicAuth
	UserAgent string
	ConnectionTimeout time.Duration
}

type Session struct {
	MaxRedirect int
	Client *http.Client
	Redirects map[string][]string
	Proxies []string
	Mutex *sync.Mutex
}

type Response struct {
	Url string
	EffectiveUrl string
	Encoding string
	Content []byte
	Cookies []*http.Cookie
	StartAt time.Time
	Elapsed time.Duration
	Header http.Header
	StatusCode int
	Redirects []string
}