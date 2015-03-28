package timesheet

import (
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	timeLayout = "2006-01-02"
	baseUrl    = "https://my.cybera.ca"
)

var (
	cookieJar        = newCookieJar()
	baseUrlCanonical = &url.URL{Scheme: "https", Host: "my.cybera.ca"}
)

func newCookieJar() *cookiejar.Jar {
	// CookieJar will never return `err`. Seriously, check the source code.
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return jar
}

func UpdateCookieJar(ctx *kingpin.ParseContext) {
	key := ""
	if splitKey := strings.Split(key, ":"); len(key) == 69 && len(splitKey) == 2 {
		sessionCookie := &http.Cookie{
			Name:   splitKey[0],
			Value:  splitKey[1],
			Domain: "my.cybera.ca",
			Path:   "/",
		}
		jsCookie := &http.Cookie{
			Name:   "has_js",
			Value:  "0",
			Domain: "my.cybera.ca",
			Path:   "/",
		}
		cookieJar.SetCookies(baseUrlCanonical, []*http.Cookie{sessionCookie, jsCookie})
	} else {
		log.Fatal(errors.New("Malformed session key"))
	}
}
