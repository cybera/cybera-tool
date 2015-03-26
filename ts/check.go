package ts

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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
