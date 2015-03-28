package timesheet

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Parse HTML from a reader
func parseResponse(res *http.Response) *goquery.Document {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

// Initialise client
func newClient() *http.Client {
	return &http.Client{Jar: cookieJar}
}
