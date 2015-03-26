package ts

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func LogHours() (err error) {
	Login()

	err = f.Parse(flag.Args()[1:])
	if err != nil {
		log.Fatal(err)
	}

	if len(*flagNotes) < 1 {
		log.Fatal(errors.New("Notes can't be empty: " + *flagNotes))
	}
	// Get login page to parse form id
	resp, err := client.Get(tsAddUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	formBuildId, exists := doc.Find("#timepunchextuni-form input[name='form_build_id']").Attr("value")
	if !exists {
		log.Fatal("Can't find form build id")
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
		"fromdate[date]": {time.Time(flagFrom).Format(dateLayout)},
		"todate[date]":   {time.Time(flagTo).Format(dateLayout)},
		"hours":          {strconv.FormatFloat(flagHours.Hours(), 'f', 2, 64)},
		"op":             {"Submit"},
		"selaccount":     {accId},
		"notes":          {*flagNotes},
		"userext":        {user},
		"destination":    {"userts/add"},
		"autorefresh":    {"0"},
		"form_build_id":  {formBuildId},
		"form_token":     {formToken},
		"form_id":        {"timepunchextuni_form"},
	}
	fmt.Printf("Logging %s hours for %s\n", tsAddRequest["hours"][0], tsAddRequest["fromdate[date]"][0])
	// Submit time entry
	resp, err = client.PostForm(tsAddUrl, tsAddRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Check for confirmation
	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}
	confirmation := doc.Find(".messages").Text()
	if len(strings.TrimSpace(confirmation)) < 1 {
		log.Fatal(errors.New("No confirmation recieved"))
	}
	// TODO: Confirm bulk entries
	return nil
}
