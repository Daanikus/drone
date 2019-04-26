package main

import (
	"os"
	"testing"
)

func TestSshExec(t *testing.T) {
	cfg := &Server{
		IP:      "39.104.53.120",
		User:    "morya",
		SshPort: 22,

		SshPrivateKey: "id_rsa",
	}
	s, err := newSshConn(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	out := os.Stdout
	s.Exec(out, "date; pwd; ls -alh", nil)
}
