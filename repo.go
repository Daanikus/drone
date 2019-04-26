package main

import (
	"net"

	"github.com/morya/drone/util"

	crypto_ssh "golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Repo struct {
	repo *git.Repository
}

func newRepo(repoURL string, repoPath string, keyFile string) (*Repo, error) {
	var exist bool
	if !util.Exists(repoPath) || !util.IsDir(repoPath) {
		util.RemoveDir(repoPath)
		exist = false
	} else {
		exist = true
	}

	var repo *git.Repository
	var err error
	if exist {
		repo, err = git.PlainOpen(repoPath)
	} else {
		// TODO 'git' 是git-server默认用户名，暂时不考虑支持其它用户名
		auth, err := ssh.NewPublicKeysFromFile("git", keyFile, "")
		if err != nil {
			return nil, err
		}
		auth.HostKeyCallback = func(hostname string, remote net.Addr, key crypto_ssh.PublicKey) error {
			// ignore host key
			return nil
		}

		repo, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:  repoURL,
			Auth: auth,
		})
	}

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
