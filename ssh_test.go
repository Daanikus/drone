package main

import (
	"os"
	"testing"

	"github.com/morya/utils/log"
)

func TestSshExec(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevelString("debug")

	var sample = "drone.json"
	if err := loadConfig(sample); err !=nil {
		t.Errorf("load config failed")
	}

	for name, serverCfg := range config.Servers {
		log.Info("before create conn")
		s, err := newSshConn(serverCfg)
		if err != nil {
			t.Error(err, "create ssh connection failed")
			return
		}

		// call cmd on server
		log.Infof("execute cmd on server [%v:%v]", name, serverCfg.IP)
		s.Exec(os.Stderr, "date; pwd; ls -alh", nil)
	}
}
