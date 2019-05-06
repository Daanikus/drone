package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type Config struct {
	Servers  map[string]*ServerConfig  `json:"servers"`
	Projects map[string]*ProjectConfig `json:"projects"`
}

var config = &Config{}

func isFileReadable(filename string) bool {
	if stat, err := os.Stat(filename); err != nil {
		if stat.IsDir() {
			return false
		}
	}
	return true
}

func expand(s string) string {
	s = os.ExpandEnv(s)
	home := os.Getenv("HOME")
	s = strings.Replace(s, "~", home, -1)
	return s
}

func isValidConfig() (err error) {
	for prjName, p := range config.Projects {
		server := p.Server
		if _, ok := config.Servers[server]; !ok {
			err = fmt.Errorf("server [%v] in project [%v] not valid", server, prjName)
			return
		}
	}
	for name, server := range config.Servers {
		if server.SshPort <= 0 {
			server.SshPort = 22
		}

		if server.Pswd == "" && server.SshPrivateKey == "" {
			err = fmt.Errorf("server [%v], must provide pswd or ssh_private_key", name)
			return
		}
		if ip := net.ParseIP(server.IP); ip == nil {
			err = fmt.Errorf("server [%v], invalid ip = [%v]", name, ip)
			return
		}

		if server.SshPrivateKey != "" {
			server.SshPrivateKey = expand(server.SshPrivateKey)
			if !isFileReadable(server.SshPrivateKey) {
				err = fmt.Errorf("server [%v] ssh_private_key = [%v] not accessable", name, server.SshPrivateKey)
				return
			}
		}
	}
	return
}

func loadConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	if err := isValidConfig(); err != nil {
		return err
	}
	return nil
}
