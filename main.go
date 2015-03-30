package main

import (
	"os"

	ts "github.com/cybera/cybera-tool/timesheet"

	"gopkg.in/alecthomas/kingpin.v1"
)

// Define the UI
var (
	app            = kingpin.New("cybera", "A tool that makes cyberans happy.")
	appCredentials = app.Flag("credentials", "Timesheet service username and password.").Short('c').PlaceHolder("USER:PASS").String()
	appKey         = app.Flag("key", "Timesheet service session key. Will be read from $CYBERA_KEY, if not provided.").PlaceHolder("SESSION_KEY").OverrideDefaultFromEnvar("CYBERA_KEY").String()

	// Log time command. Not to confuse with `log.Fatal()`.
	l            = app.Command("log", "Log a timesheet entry(s).")
	lTime        = l.Flag("time", "Duration to log.").Short('t').Default("7h").Duration()
	lAccount     = l.Arg("account", "Account under which time will be logged.").Required().String()
	lDescription = l.Arg("description", "Description of work.").Required().String()
	lDates       = l.Arg("dates", "Date or date range to log this work at. TO BE IMPLEMENTED.").Default(ts.Today).Strings()
)

func main() {
	app.Version("0.0.1")
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	if len(*appKey) != 0 {
		ts.UpdateCookieJar(*appKey)
	} else {
		if len(*appCredentials) != 0 {
			ts.Auth(*appCredentials)
		} else {
			println("Can't do much without either session key or user credentials.")
			os.Exit(1)
		}
	}
	switch command {
	case l.FullCommand():
		ts.Log(*lDates, *lTime, *lAccount, *lDescription)
	}
}
