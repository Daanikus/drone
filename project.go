package main

type Project struct {
	GitURL string `json:"git_url"`
	Server string

	PreTasks   []Task `json:"pre_task"`
	BuildTasks []Task `json:"task"`
	PostTasks  []Task `json:"post_task"`
}
