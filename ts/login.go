package ts

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseUrl  = "https://my.cybera.ca"
	loginUrl = baseUrl + "/login"
	otUrl    = baseUrl + "/userto"
)

// Define network client
var jar *cookiejar.Jar
var client http.Client
var user = os.Getenv("TS_USER")
var pass = os.Getenv("TS_PWD")

func Login() {
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
		log.Fatal("Can't find form build id")
	}

	// Create a form
	loginRequest := url.Values{
		"name":          {user},
		"pass":          {pass},
		"form_build_id": {formBuildId},
		"form_id":       {"user_login"},
		"op":            {"Log in"},
	}

	// Log in
	resp, err = client.PostForm(loginUrl, loginRequest)
	if err != nil {
		log.Fatal(err)
	}
}
