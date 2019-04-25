package main

type Repo struct {
}

func newRepo(path string) *Repo {
	return &Repo{}
}

func (r *Repo) HasUpdate() bool {
	return false
}
