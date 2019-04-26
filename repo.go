package main

import (
	git "gopkg.in/src-d/go-git.v4"
	git_memory_storage "gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/morya/utils/log"
)

type Repo struct {
	repo *git.Repository
}

func newRepo(path string) *Repo {
	repoStorage := git_memory_storage.NewStorage()
	repo, err := git.Init(repoStorage, nil)
	if err != nil {
		log.InfoError(err)
		return nil
	}
	return &Repo{repo: repo}
}

func (r *Repo) HasUpdate() bool {
	return false
}

func (r *Repo) Pull() bool {
	return false
}
