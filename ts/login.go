package ts

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseUrl    = "https://my.cybera.ca"
	loginUrl   = baseUrl + "/login"
	tsUrl      = baseUrl + "/userts"
	tsAddUrl   = tsUrl + "/add"
	otUrl      = baseUrl + "/userto"
	user       = "valiushko"
	pass       = "valiushko"
	dateLayout = "2006-01-02"
)

// Define network client
var jar *cookiejar.Jar
var client http.Client

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
