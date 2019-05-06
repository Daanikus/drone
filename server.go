package main

type ServerConfig struct {
	IP            string
	User          string
	Pswd          string `json:"pswd"`
	SshPrivateKey string `json:"ssh_private_key"`
	SshPort       int    `json:"ssh_port"`
}
