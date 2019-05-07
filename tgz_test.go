package main

import (
	"testing"

	"github.com/morya/utils/log"
	"gopkg.in/src-d/go-git.v4"
)

func TestTarGzipClone(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevelString("debug")
	tar := &TarGzip{}
	repo, err := git.PlainOpen(`D:\gocode\mailfile`)
	if err != nil {
		log.Info(err)
		return
	}
	src, _ := repo.Worktree()
	tar.Clone(src.Filesystem)
}
