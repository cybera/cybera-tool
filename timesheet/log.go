package timesheet

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"illotum/cybera-tool/session"

	"github.com/PuerkitoBio/goquery"
)

const (
	tsUrl    = baseUri + "/userts"
	tsAddUrl = tsUrl + "/add"
)

func Log() {
	Auth("valiushko", "valiushko")
	if err := f.Parse(flag.Args()[1:]); err != nil {
		log.Fatal(err)
	}

	if len(*flagDesc) < 1 {
		log.Fatal(errors.New("Notes can't be empty"))
	}
	// Get login page to parse form id
	fmt.Println(session.Key.Value)
	res, keyUpdated := session.Get(tsAddUrl)
	if keyUpdated {
		fmt.Println("KEY UPDATED: ", session.Key.Value)
	}
	doc := session.Parse(res)

	formBuildId, exists := doc.Find("#timepunchextuni-form input[name='form_build_id']").Attr("value")
	if !exists {
		log.Fatal("Can't find timesheet form build id")
	}

	formToken, exists := doc.Find("#edit-timepunchextuni-form-form-token").Attr("value")
	if !exists {
		log.Fatal("Can't find form token")
	}

	accounts := make(map[string]string)
	doc.Find("#edit-selaccount").ChildrenFiltered("option").Each(
		func(i int, opt *goquery.Selection) {
			accounts[opt.Text()], _ = opt.Attr("value")
		})

	accId := accounts[*flagAccount]
	if accId == "" {
		log.Fatal(errors.New("Account not found"))
	}

	// Fill the form
	tsAddRequest := url.Values{
		"fromdate[date]": {time.Time(flagAt).Format(dateLayout)},
		"todate[date]":   {time.Time(flagAt).Format(dateLayout)},
		"hours":          {strconv.FormatFloat(flagHours.Hours(), 'f', 2, 64)},
		"op":             {"Submit"},
		"selaccount":     {accId},
		"notes":          {*flagDesc},
		"userext":        {},
		"destination":    {"userts/add"},
		"autorefresh":    {"0"},
		"form_build_id":  {formBuildId},
		"form_token":     {formToken},
		"form_id":        {"timepunchextuni_form"},
	}
	fmt.Printf("Logging %s hours for %s\n", tsAddRequest["hours"][0], tsAddRequest["fromdate[date]"][0])
	// Submit time entry
	res, keyUpdated = session.Post(tsAddUrl, tsAddRequest)
	if keyUpdated {
		fmt.Println("KEY UPDATED: ", session.Key.Value)
	}
	doc = session.Parse(res)

	confirmation := doc.Find(".messages").Text()
	if len(strings.TrimSpace(confirmation)) < 1 {
		log.Fatal(errors.New("No confirmation recieved"))
	}
}
