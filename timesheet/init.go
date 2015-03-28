package timesheet

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

const (
	dateLayout = "2006-01-02"
)

var (
	// Flags
	f                    = flag.NewFlagSet("ts", flag.PanicOnError)
	flagHours            = f.Duration("time", time.Duration(7)*time.Hour, "Duration to log")
	flagAccount          = f.String("acc", "DevOps", "Account to log under")
	flagDesc             = f.String("desc", "", "Description of work")
	flagUser             = f.String("user", "", "Username for login")
	flagPass             = f.String("pass", "", "Password for login")
	flagAt      timeFlag = timeFlag(time.Now())
	cookieJar   *cookiejar.Jar
	baseUrl     = url.URL{
		Scheme: "https",
		Host:   "my.cybera.ca",
	}
)

type timeFlag time.Time

func (t *timeFlag) Set(s string) error {
	v, err := time.Parse(dateLayout, s)
	*t = timeFlag(v)
	return err
}
func (t *timeFlag) String() string { return (*time.Time)(t).String() }

func init() {
	// Setup a date flag
	f.Var(&flagAt, "at", "Submit time for this date")

	// Setup session
	cookieJar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	if envKey := os.Getenv("TS_KEY"); len(envKey) != 0 {
		if splitKey := strings.Split(os.Getenv("TS_KEY"), ":"); len(envKey) == 69 && len(splitKey) == 2 {
			key := &http.Cookie{
				Name:   splitKey[0],
				Value:  splitKey[1],
				Domain: "my.cybera.ca",
				Path:   "/",
			}
			cookieJar.SetCookies(
				&url.URL{Scheme: "https", Host: "my.cybera.ca"},
				[]*http.Cookie{key},
			)
		} else {
			log.Fatal(errors.New("Malformed session key"))
		}
	}
}
