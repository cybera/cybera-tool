package ts

import (
	"log"
	"net/url"
)

const (
	baseUri  = "https://my.cybera.ca"
	loginUri = baseUri + "/login"
	otUrl    = baseUri + "/userto"
)

func Login() {
	// Get login page to parse form id
	doc := getParsed(loginUri)
	formBuildId, exists := doc.Find("#user-login div input[type='hidden']").Attr("id")
	if !exists {
		log.Fatal("Can't find login form build id")
	}

	// Create a form
	form := url.Values{
		"name":          {user},
		"pass":          {pass},
		"form_build_id": {formBuildId},
		"form_id":       {"user_login"},
		"op":            {"Log in"},
	}

	// Log in
	postParsed(loginUri, form)

	// Store cookies
}
