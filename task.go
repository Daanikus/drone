package main

import "github.com/morya/utils/log"

type Task struct {
	Command      string            `json:"command"`
	Env          map[string]string `json:"env"`
	IsReplaceEnv bool              `json:"is_replace_env"`
	Log          *log.Logger       `json:"-"`
}
