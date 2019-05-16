package main

import "time"

type Event struct {
	Name      string
	PrjConfig *ProjectConfig
	// Commit  string // if specified, trigger special commit from git project
	Created time.Time
}
