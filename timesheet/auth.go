package timesheet

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	authUrl = baseUrl + "/login"
)

func Auth(credentials string) {
	var user, pass string
	if splitCredentials := strings.Split(credentials, ":"); len(splitCredentials) == 2 {
		user, pass = splitCredentials[0], splitCredentials[1]
	} else {
		log.Fatal(errors.New("Malformed credentials"))
	}

	c := newClient()
	// Get login page to parse the form id
	res, _ := c.Get(authUrl)
	doc := parseResponse(res)
	formBuildId, exists := doc.Find("#user-login div input[name='form_build_id']").Attr("id")
	if !exists {
		log.Fatal("Can't find login form build id")
	}

	// Fill the form
	form := url.Values{
		"name":          {user},
		"pass":          {pass},
		"form_build_id": {formBuildId},
		"form_id":       {"user_login"},
		"op":            {"Log in"},
	}

	// Log in
	res, _ = c.PostForm(authUrl, form)

	// Feedback
	doc = parseResponse(res)
	msgs := doc.Find(".error")
	for i := range msgs.Nodes {
		m := msgs.Eq(i)
		println(strings.TrimSpace(m.Text()))
	}
	if res.StatusCode == 200 {
		cookies := c.Jar.Cookies(baseUrlCanonical)
		var key *http.Cookie
		for _, v := range cookies {
			if v.Name[:4] == "SESS" {
				key = v
			}
		}
		fmt.Printf("New session key:\n%s:%s\n", key.Name, key.Value)
	} else {
		fmt.Printf("Login failed: %v\n", res.Status)
	}
}
