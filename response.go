package requests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/clbanning/x2j"
	"github.com/sloonz/go-iconv"
	"io"
	"strings"
)

func (this *Response) Text() string {
	return string(this.Content)
}

func (this *Response) Json(obj interface{}) error {
	err := json.Unmarshal(this.Content, obj)
	return err
}

func (this *Response) Xml() (res map[string]interface{}, err error) {
	return x2j.DocToMap(this.Text())
}

func (this *Response) Feed() (*FeedChannel, error) {
	xmlDecoder := xml.NewDecoder(bytes.NewReader(this.Content))
	xmlDecoder.CharsetReader = UTF8Reader
	var rss struct {
		Channel FeedChannel `xml:"channel"`
	}
	err := xmlDecoder.Decode(&rss)
	if err != nil {
		return nil, err
	}
	return &rss.Channel, nil
}

func UTF8Reader(charset string, r io.Reader) (io.Reader, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}
	converted, err := iconv.Conv(buf.String(), "UTF-8", charset)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(converted), nil
}
