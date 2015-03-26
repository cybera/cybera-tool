package ts

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl    = "https://my.cybera.ca"
	loginUrl   = baseUrl + "/login"
	tsUrl      = baseUrl + "/userts"
	tsAddUrl   = tsUrl + "/add"
	otUrl      = baseUrl + "/userto"
	user       = "valiushko"
	pass       = "valiushko"
	dateLayout = "2006-01-02"
)

var jar *cookiejar.Jar
var client http.Client

func Init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{Jar: jar}

	// Get login page to parse form id
	resp, err := client.Get(loginUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	formBuildId, exists := doc.Find("#user-login input[type='hidden']").Attr("id")
	if !exists {
		log.Fatal("Can't find form build id")
	}

	// Create a form
	loginRequest := url.Values{
		"name":          {user},
		"pass":          {pass},
		"form_build_id": {formBuildId},
		"form_id":       {"user_login"},
		"op":            {"Log in"},
	}

	// Log in
	resp, err = client.PostForm(loginUrl, loginRequest)
	if err != nil {
		log.Fatal(err)
	}
}

func OvertimeReport() {
	resp, err := client.Get(otUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", doc.Text())
}

func CheckTimesheet() (float64, float64, []time.Time) {
	reRequired := regexp.MustCompile("required = ([0-9.]+)")
	reEntered := regexp.MustCompile("entered = ([0-9.]+)")
	reMissing := regexp.MustCompile("dates: ([0-9 ,-]+)")

	resp, err := client.Get(tsUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	report := doc.Find(".view-vw-TimeslicesMaxi .view-header p").Text()
	hoursRequiredString := reRequired.FindStringSubmatch(report)[1]
	hoursEnteredString := reEntered.FindStringSubmatch(report)[1]
	missingDatesStrings := strings.Split(reMissing.FindStringSubmatch(report)[1], ", ")

	hoursRequired, err := strconv.ParseFloat(hoursRequiredString, 64)
	if err != nil {
		log.Fatal(err)
	}
	hoursEntered, err := strconv.ParseFloat(hoursEnteredString, 64)
	if err != nil {
		log.Fatal(err)
	}
	missingDates := make([]time.Time, len(missingDatesStrings))
	for i, v := range missingDatesStrings {
		missingDates[i], err = time.Parse(dateLayout, strings.TrimSpace(v))
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%f ", hoursRequired)
	fmt.Printf("%f\n", hoursEntered)
	fmt.Printf("%s", missingDates)
	return hoursRequired, hoursEntered, missingDates
}

func LogHours(from, to time.Time, hours float64, acc, notes string) (err error) {
	if len(notes) < 1 {
		return errors.New("Notes can't be empty")
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

	accId := accounts[acc]
	if accId == "" {
		return errors.New("Account not found")
	}

	// Fill the form
	tsAddRequest := url.Values{
		"fromdate[date]": {from.Format(dateLayout)},
		"todate[date]":   {to.Format(dateLayout)},
		"hours":          {strconv.FormatFloat(hours, 'f', 2, 64)},
		"op":             {"Submit"},
		"selaccount":     {accId},
		"notes":          {notes},
		"userext":        {user},
		"destination":    {"userts/add"},
		"autorefresh":    {"0"},
		"form_build_id":  {formBuildId},
		"form_token":     {formToken},
		"form_id":        {"timepunchextuni_form"},
	}

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
		return errors.New("No confirmation recieved")
	}
	// TODO: Confirm bulk entries
	return nil
}
