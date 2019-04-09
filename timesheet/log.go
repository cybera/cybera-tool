package timesheet

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	logUrl = baseUrl + "/index.php?newent=y"
)

var Today = time.Now().Format(timeLayout)

func parseAccounts(doc *goquery.Document) map[string]string {
	accounts := make(map[string]string)
	doc.Find("select").ChildrenFiltered("option").Each(
		func(i int, opt *goquery.Selection) {
			accounts[opt.Text()], _ = opt.Attr("value")
		})
	return accounts
}

func Log(dates []string, time time.Duration, account, descr string, noop bool) {
	c := newClient()

	// Get login page to parse the form id
	res, _ := c.Get(logUrl)
	doc := parseResponse(res)

	//Deal with csrf
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

	// The form expects an account ID, not name
	accounts := parseAccounts(doc)
	log.Println(accounts)
	accId := accounts[account]
	if accId == "" {
		log.Fatal(errors.New("Account not found."))
	}

	// Fill the form
	tsAddRequest := url.Values{
		"entry_date":        {dates[0]}, // TODO: UNSAFE
		"entry_todate":      {dates[0]}, // TODO: Implement ranges
		"entry_hours":       {strconv.FormatFloat(time.Hours(), 'f', 2, 64)},
		"go":                {"Create Timesheet entry"},
		"account_id":        {accId},
		"entry_description": {descr},
		"add":               {"y"},
		"ctsrtx":            {csrf_value},
	}
	fmt.Printf("Logging %s hours for %s\n", tsAddRequest["entry_hours"][0], tsAddRequest["entry_date"][0])

	if noop {
		fmt.Printf("Would have POST'd: %+v\n", tsAddRequest)
	} else {
		// Submit time entry
		res, _ = c.PostForm(logUrl, tsAddRequest)

		doc = parseResponse(res)
		confirmation := doc.Find("html").Text()
		log.Println(confirmation)
		if len(strings.TrimSpace(confirmation)) < 1 {
			log.Fatal(errors.New("No confirmation recieved"))
		}
	}
}
