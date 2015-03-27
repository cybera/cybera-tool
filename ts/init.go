package ts

import (
	"net/http"
	"net/http/cookiejar"
	"os"

	"golang.org/x/net/publicsuffix"
)

var (
	cookieJar *cookiejar.Jar
	client    http.Client
	user      = os.Getenv("TS_USER")
	pass      = os.Getenv("TS_PWD")
)

// Init global HTTP client and Cookie storage
func init() {
	cookieJar, _ = cookiejar.New(
		&cookiejar.Options{PublicSuffixList: publicsuffix.List},
	)

	client = http.Client{Jar: cookieJar}
}
