package main

import (
	"os"
	"testing"
)

func TestSshExec(t *testing.T) {
	cfg := &Server{
		IP:      "192.168.2.156",
		User:    "morya",
		SshPort: 22,

		SshPrivateKey: "/Users/morya/.ssh/id_rsa",
	}
	s, err := newSshConn(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	out := os.Stdout
	s.Exec(out, "date", nil)
}
