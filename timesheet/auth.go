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
	authUrl = baseUrl
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
	csrf_value := ""

	// CSRF - is stored in the cookie
	if res.StatusCode == 200 {
		cookies := c.Jar.Cookies(baseUrlCanonical)
		for _, v := range cookies {
			if v.Name == "ctsrtx" {
				csrf_value = v.Value
			}
		}
	}

	// Fill the form
	form := url.Values{
		"auth_user": {user},
		"auth_pw":   {pass},
		//"form_build_id": {formBuildId},
		//"form_id":    {"user_login"},
		"auth_login": {"Log in"},
		"ctsrtx":     {csrf_value},
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
			if v.Name == "PHPSESSID" {
				key = v
			}
		}
		if key != nil {
			fmt.Printf("New session key:\n%s:%s\n", key.Name, key.Value)
		} else {
			fmt.Printf("Failed to find session key")
		}
	} else {
		fmt.Printf("Login failed: %v\n", res.Status)
	}
}
