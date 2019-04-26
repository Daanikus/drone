package main

import (
	"net"

	"gopkg.in/src-d/go-git.v4/storage/filesystem"

	"gopkg.in/src-d/go-billy.v4"

	crypto_ssh "golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Repo struct {
	repo *git.Repository
}

func newRepo(repoURL string, repoPath string, keyFile string) (*Repo, error) {
	var wt, dot billy.Filesystem

	wt = osfs.New(repoPath)
	dot, _ = wt.Chroot(".git")

	// TODO 'git' 是git-server默认用户名，暂时不考虑支持其它用户名
	auth, err := ssh.NewPublicKeysFromFile("git", keyFile, "")
	if err != nil {
		return nil, err
	}
	auth.HostKeyCallback = func(hostname string, remote net.Addr, key crypto_ssh.PublicKey) error {
		// ignore host key
		return nil
	}

	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	repo, err := git.Clone(storage, wt, &git.CloneOptions{
		URL:  repoURL,
		Auth: auth,
	})
	if err != nil {
		return nil, err
	}
	return &Repo{repo: repo}, nil
}

func (r *Repo) HasUpdate() bool {
	return true
}

func (r *Repo) Pull() bool {
	return false
}
