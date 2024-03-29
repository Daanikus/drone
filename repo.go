package main

import (
	"net"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/morya/drone/util"
	"github.com/morya/utils/log"
	"github.com/pkg/errors"
	crypto_ssh "golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Repo struct {
	keyFile    string
	repo       *git.Repository
	lastCommit *object.Commit
}

func getRepoLastCommit(r *git.Repository) (*object.Commit, error) {
	logs, _ := r.Log(&git.LogOptions{})
	var lastCommit *object.Commit
	logs.ForEach(func(commit *object.Commit) error {
		if lastCommit == nil {
			lastCommit = commit
			return errors.New("done")
		}
		// log.Debugf("commit hash = %v", commit.Hash.String())
		return nil
	})

	return lastCommit, nil
}

func getAuth(keyFile string) (transport.AuthMethod, error) {
	// TODO 'git' 是git-server默认用户名
	//   暂不支持其它用户名
	auth, err := ssh.NewPublicKeysFromFile("git", keyFile, "")
	if err != nil {
		return nil, errors.Wrapf(err, "read git key failed, key=%v", keyFile)
	}
	auth.HostKeyCallback = func(hostname string, remote net.Addr, key crypto_ssh.PublicKey) error {
		// ignore host key
		return nil
	}

	return auth, nil
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

		auth, err := getAuth(keyFile)
		if err != nil {
			return nil, errors.Wrap(err, "create auth failed")
		}

		repo, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:  repoURL,
			Auth: auth,
		})
	}

	if err != nil {
		return nil, errors.Wrapf(err, "open git repo failed, path=%v", repoPath)
	}
	if repo == nil {
		return nil, errors.New("repo is nil")
	}
	lastCommit, err := getRepoLastCommit(repo)
	log.Debugf("last commit is %v, err = %v", lastCommit.Hash.String(), err)
	return &Repo{keyFile: keyFile, repo: repo, lastCommit: lastCommit}, nil
}

func (r *Repo) HasUpdate() (bool, error) {
	auth, err := getAuth(r.keyFile)
	if err != nil {
		return false, err
	}

	// r.repo.Fetch(&git.FetchOptions{RemoteName: "origin", Auth: auth})

	tree, err := r.repo.Worktree()
	if err != nil {
		return false, err
	}

	err = tree.Pull(&git.PullOptions{RemoteName: "origin", Auth: auth})
	switch err {
	case git.NoErrAlreadyUpToDate:
		return false, nil

	case nil:
		return true, nil
	default:

	}
	return false, err
}

// return nil on success
// return error on fail
func (r *Repo) PackZipFile(filename string) (err error) {
	root, _ := r.repo.Worktree()
	gzip := &TarGzip{}
	gzip.Clone(root.Filesystem)
	return
}
