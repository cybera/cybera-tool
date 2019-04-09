package main

import (
	"os"

	ts "github.com/cybera/cybera-tool/timesheet"

	"gopkg.in/alecthomas/kingpin.v1"
)

// Define the UI
var (
	app    = kingpin.New("cybera", "Tool that makes cyberans happy.")
	appKey = app.Flag("key", "Timesheet service session key. Will be read from $CYBERA_KEY, if not provided.").PlaceHolder("SESSION_KEY").OverrideDefaultFromEnvar("CYBERA_KEY").String()

	// authenticate subcommand
	a            = app.Command("authenticate", "Authenticate and get a session key.")
	aCredentials = a.Flag("credentials", "Timesheet service username and password.").PlaceHolder("USER:PASS").OverrideDefaultFromEnvar("CYBERA_CREDS").Short('c').String()

	// Log time command. Not to confuse with `log.Fatal()`.
	l            = app.Command("log", "Log a timesheet entry(s).")
	lTime        = l.Flag("time", "Duration to log.").Short('t').Default("7h").Duration()
	lNoop        = l.Flag("noop", "no-op mode. Don't actually post a time.").Short('n').Bool()
	lAccount     = l.Arg("account", "Account under which time will be logged.").Required().String()
	lDescription = l.Arg("description", "Description of work.").Required().String()
	lDates       = l.Arg("dates", "Date or date range to log this work at. TO BE IMPLEMENTED.").Default(ts.Today).Strings()
)

func main() {
	app.Version("0.0.2")
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	if len(*appKey) != 0 {
		ts.UpdateCookieJar(*appKey)
	} else {
		if len(*aCredentials) == 0 {
			println("Can't do much without either session key or user credentials.")
			os.Exit(1)
		}
	}
	switch command {
	case a.FullCommand():
		if len(*aCredentials) != 0 {
			ts.Auth(*aCredentials)
		}
	case l.FullCommand():
		ts.Log(*lDates, *lTime, *lAccount, *lDescription, *lNoop)
	}
}
