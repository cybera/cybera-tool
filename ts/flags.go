package ts

import (
	"flag"
	"time"
)

const (
	dateLayout = "2006-01-02"
)

// Define command-line flags
var (
	f           = flag.NewFlagSet("timesheet", flag.PanicOnError)
	flagHours   = f.Duration("hours", time.Duration(7)*time.Hour, "Duration to log")
	flagAccount = f.String("account", "DevOps", "Account to log under")
	flagNotes   = f.String("notes", "", "Description of work")
)

type timeFlag time.Time

func (t *timeFlag) Set(s string) error {
	v, err := time.Parse(dateLayout, s)
	*t = timeFlag(v)
	return err
}
func (t *timeFlag) String() string { return (*time.Time)(t).String() }

var flagFrom timeFlag = timeFlag(time.Now())
var flagTo timeFlag = flagFrom

func init() {
	f.Var(&flagFrom, "from", "Start date")
	f.Var(&flagTo, "to", "End date")
}
