package ts

import (
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/franela/goreq"
)

// Get, check for errors and parse
func getParsed(uri string) *goquery.Document {
	res, err := goreq.Request{
		CookieJar:   cookieJar,
		Compression: goreq.Gzip(),
		Uri:         uri,
	}.Do()
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// Post, check for errors and parse
func postParsed(uri string, form url.Values) *goquery.Document {
	c := http.Client{Jar: cookieJar}
	res, err := c.PostForm(uri, form)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}
