package ts

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	baseUrl  = "https://my.cybera.ca"
	loginUrl = baseUrl + "/login"
	tsUrl    = baseUrl + "/userts"
	otUrl    = baseUrl + "/userto"
	user     = "valiushko"
	pass     = "valiushko"
)

var jar *cookiejar.Jar
var client http.Client

func Init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{Jar: jar}

	// Get login page to parse form id
	resp, err := client.Get(loginUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	formBuildId, exists := doc.Find("#user-login input[type='hidden']").Attr("id")
	if !exists {
		log.Fatal("Can't find login form id")
	}

	// Create a form
	loginRequest := url.Values{
		"name":          {user},
		"pass":          {pass},
		"form_build_id": {formBuildId},
		"form_id":       {"user_login"},
		"op":            {"Log+in"},
	}

	// Log in
	resp, err = client.PostForm(loginUrl, loginRequest)
	if err != nil {
		log.Fatal(err)
	}
}

func OvertimeReport() {
	resp, err := client.Get(otUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	fmt.Printf("%s", doc.Text())
}

func CheckTimesheet() {
	resp, err := client.Get(tsUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	fmt.Printf("%s", doc.Text())
}
