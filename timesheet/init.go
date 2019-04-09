package timesheet

import (
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const (
	timeLayout = "2006-01-02"
	baseUrl    = "https://timesheet.cybera.ca"
)

var (
	cookieJar        = newCookieJar()
	baseUrlCanonical = &url.URL{Scheme: "https", Host: "timesheet.cybera.ca"}
)

func newCookieJar() *cookiejar.Jar {
	// CookieJar will never return `err`. Seriously, check the source code.
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return jar
}

func UpdateCookieJar(key string) {
	if splitKey := strings.Split(key, ":"); len(key) == 36 && len(splitKey) == 2 {
		sessionCookie := &http.Cookie{
			Name:   splitKey[0],
			Value:  splitKey[1],
			Domain: "timesheet.cybera.ca",
			Path:   "/",
		}
		csrfCookie := &http.Cookie{}
		cookieJar.SetCookies(baseUrlCanonical, []*http.Cookie{sessionCookie, csrfCookie})
	} else {
		log.Fatal(errors.New("Malformed session key"))
	}
}
