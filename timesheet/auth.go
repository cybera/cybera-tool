package timesheet

import (
	"fmt"
	"illotum/cybera-tool/session"
	"log"
	"net/url"
	"strings"
)

const (
	baseUri  = "https://my.cybera.ca"
	loginUri = baseUri + "/login"
	otUrl    = baseUri + "/userto"
)

func Auth(user, pass string) {
	// Get login page to parse form id
	fmt.Println(session.Key)
	res, keyUpdated := session.Get(loginUri)
	fmt.Println(session.Key)
	// Do not process anything unless we got a new session key
	doc := session.Parse(res)
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
	fmt.Println(session.Key)
	res, keyUpdated = session.Post(loginUri, form)
	if keyUpdated && res.StatusCode == 200 {
		fmt.Printf("To avoid repeated logins, please export:\nTS_KEY=%s:%s\n", session.Key.Name, session.Key.Value)
	} else {
		fmt.Println(session.Key)
		fmt.Printf("Login failed: %v\n", res.Status)
		// Not updated key means that session hasn't changed, i.e. login failed
		doc = session.Parse(res)
		msg := doc.Find(".error").First()
		if msg != nil {
			fmt.Printf("%s\n", strings.TrimSpace(msg.Text()))
		}
	}
}

func Auth2(user, pass string) {
	c := session.NewClient()
	// Get login page to parse form id
	fmt.Println(session.Key)
	res, _ := c.Get(loginUri)
	fmt.Println(session.Key)
	// Do not process anything unless we got a new session key
	doc := session.Parse(res)
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
	fmt.Println(session.Key)
	res, _ = c.PostForm(loginUri, form)
	fmt.Println(session.Key)
	if res.StatusCode == 200 {
		// fmt.Printf("To avoid repeated logins, please export:\nTS_KEY=%s:%s\n", session.Key.Name, session.Key.Value)
	} else {
		fmt.Println(session.Key)
		fmt.Printf("Login failed: %v\n", res.Status)
		// Not updated key means that session hasn't changed, i.e. login failed
		doc = session.Parse(res)
		msg := doc.Find(".error").First()
		if msg != nil {
			fmt.Printf("%s\n", strings.TrimSpace(msg.Text()))
		}
	}
}
