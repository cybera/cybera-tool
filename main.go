package main

import (
	"illotum/cybera/ts"
	"time"
)

func main() {
	ts.Init()
	ts.LogHours(time.Now(), time.Now(), 7, "DevOps", "Some project")
}
