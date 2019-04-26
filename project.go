package main

type Project struct {
	Server    string
	GitURL    string `json:"git_url"`
	GitKey    string `json:"git_key"`
	LocalPath string `json:"local_path"`

	PreTasks   []Task `json:"pre_task"`
	BuildTasks []Task `json:"task"`
	PostTasks  []Task `json:"post_task"`
}
