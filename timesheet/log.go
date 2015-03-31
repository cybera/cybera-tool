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
	logUrl = baseUrl + "/userts/add"
)

var Today = time.Now().Format(timeLayout)

func parseAccounts(doc *goquery.Document) map[string]string {
	accounts := make(map[string]string)
	doc.Find("#edit-selaccount").ChildrenFiltered("option").Each(
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

	formBuildId, exists := doc.Find("#timepunchextuni-form input[name='form_build_id']").Attr("value")
	if !exists {
		log.Fatal("Can't find timesheet form build id")
	}
	formToken, exists := doc.Find("#edit-timepunchextuni-form-form-token").Attr("value")
	if !exists {
		log.Fatal("Can't find form token")
	}

	// The form expects an account ID, not name
	accounts := parseAccounts(doc)
	accId := accounts[account]
	if accId == "" {
		log.Fatal(errors.New("Account not found."))
	}

	// Fill the form
	tsAddRequest := url.Values{
		"fromdate[date]": {dates[0]}, // TODO: UNSAFE
		"todate[date]":   {dates[0]}, // TODO: Implement ranges
		"hours":          {strconv.FormatFloat(time.Hours(), 'f', 2, 64)},
		"op":             {"Submit"},
		"selaccount":     {accId},
		"notes":          {descr},
		"userext":        {},
		"destination":    {"userts/add"},
		"autorefresh":    {"0"},
		"form_build_id":  {formBuildId},
		"form_token":     {formToken},
		"form_id":        {"timepunchextuni_form"},
	}
	fmt.Printf("Logging %s hours for %s\n", tsAddRequest["hours"][0], tsAddRequest["fromdate[date]"][0])

	if noop {
		fmt.Printf("Would have POST'd: %+v\n", tsAddRequest)
	} else {
		// Submit time entry
		res, _ = c.PostForm(logUrl, tsAddRequest)

		doc = parseResponse(res)
		confirmation := doc.Find(".messages").Text()
		if len(strings.TrimSpace(confirmation)) < 1 {
			log.Fatal(errors.New("No confirmation recieved"))
		}
	}
}
