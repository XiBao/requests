package requests

import (
	"time"
)

// RSS & ATOM FEED
type FeedChannel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate FeedDate   `xml:"lastBuildDate"`
	Item          []FeedItem `xml:"item"`
}

type FeedItemEnclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type FeedItem struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Comments    string        `xml:"comments"`
	PubDate     FeedDate      `xml:"pubDate"`
	GUID        string        `xml:"guid"`
	Category    []string      `xml:"category"`
	Enclosure   FeedItemEnclosure `xml:"enclosure"`
	Description string        `xml:"description"`
	Content     string        `xml:"content"`
}

type FeedDate string

func (self FeedDate) Parse() (time.Time, error) {
	// Wordpress format
	t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", string(self)) 
	if err != nil {
		t, err = time.Parse(time.RFC822, string(self)) // RSS 2.0 spec
	}
	return t, err
}

func (self FeedDate) Format(format string) (string, error) {
	t, err := self.Parse()
	if err != nil {
		return "", err
	}
	return t.Format(format), nil
}

func (self FeedDate) MustFormat(format string) string {
	s, err := self.Format(format)
	if err != nil {
		return err.Error()
	}
	return s
}