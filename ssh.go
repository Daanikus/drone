package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/morya/utils/log"
	"golang.org/x/crypto/ssh"
)

type SshConn struct {
	conn *ssh.Client
}

func loadKey(keyPath string) ssh.AuthMethod {
	f, err := os.Open(keyPath)
	if err != nil {
		log.InfoError(err)
		return nil
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	s, _ := ssh.ParsePrivateKey(data)
	if s == nil {
		log.Infof("parse private key return nil")
		return nil
	}
	return ssh.PublicKeys(s)
}

func newSshConn(serverCfg *ServerConfig) (*SshConn, error) {
	AuthMethods := make([]ssh.AuthMethod, 0)
	if serverCfg.Pswd != "" {
		log.Debugf("will try password auth")
		AuthMethods = append(AuthMethods, ssh.Password(serverCfg.Pswd))
	}
	if serverCfg.SshPrivateKey != "" {
		log.Debugf("will try PrivateKey auth %v", serverCfg.SshPrivateKey)
		if keyAuth := loadKey(serverCfg.SshPrivateKey); keyAuth != nil {
			AuthMethods = append(AuthMethods, keyAuth)
		}
	}

	sshCfg := &ssh.ClientConfig{
		User:            serverCfg.User,
		Auth:            AuthMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	var remoteAddr = fmt.Sprintf("%v:%v", serverCfg.IP, serverCfg.SshPort)
	log.Debugf("dial remote %v", remoteAddr)
	conn, err := ssh.Dial("tcp", remoteAddr, sshCfg)
	if err != nil {
		return nil, err
	}
	return &SshConn{conn: conn}, nil
}

func (p *SshConn) Exec(output io.Writer, cmd string, env map[string]string) {
	session, err := p.conn.NewSession()
	if err != nil {
		log.InfoError(err)
		return
	}
	if env != nil {
		for k, v := range env {
			session.Setenv(k, v)
		}
	}
	data, err := session.CombinedOutput(cmd)
	if err != nil {
		log.InfoError(err)
		return
	}
	output.Write(data)
	return
}
