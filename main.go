package main

import (
	"flag"
	"illotum/cybera/ts"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "log":
		ts.LogHours()
	}
}
