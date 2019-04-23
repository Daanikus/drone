package main

type Project struct {
	GitURL string

	PreTasks   []Task
	BuildTasks []Task
	PostTasks  []Task
}
