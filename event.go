package main

import "time"

type Event struct {
	Project string
	Commit  string // if specified, trigger special commit from git project
	Created time.Time
}
