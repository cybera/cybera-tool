package ts

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func OvertimeReport() {
	resp, err := client.Get(otUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", doc.Text())
}
