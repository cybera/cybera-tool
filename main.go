package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2-unstable"
)

// Define the UI
var (
	app     = kingpin.New("cybera", "A tool that makes cyberans happy.")
	appUser = app.Flag("user", "Timesheet service username.").String()
	appPass = app.Flag("pass", "Timesheet service password.").String()
	appKey  = app.Flag("key", "Timesheet service session key. Will be read from $CYBERA_KEY, if not provided.").OverrideDefaultFromEnvar("CYBERA_KEY").String()

	// Log time command. Not to confuse with `log.Fatal()`.
	l            = app.Command("log", "Log a timesheet entry(s).")
	lTime        = l.Flag("time", "Duration to log.").Short('t').Default("7h").Duration()
	lProject     = l.Arg("project", "Project under which time will be logged.").Required().String()
	lDescription = l.Arg("description", "Description of work.").Required().String()
	lDates       = l.Arg("dates", "Date or date range to log this work at.").Strings()
)

func main() {
	app.Version("0.0.1")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case l.FullCommand():
		// timesheet.Log()
	}
}
