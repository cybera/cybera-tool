package timesheet

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Parse HTML from a reader
func parse(res *http.Response) *goquery.Document {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

// Initialise client
func New() *http.Client {
	return &http.Client{Jar: Jar}
}
