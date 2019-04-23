package main

import "github.com/morya/utils/log"

type Task struct {
	Name         string
	Command      string
	Env          map[string]string
	IsReplaceEnv bool
	Log          *log.Logger
}
