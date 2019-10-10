package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func publicKeyAuthFunc(pemBytes, keyPassword []byte) ssh.AuthMethod {
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKeyWithPassphrase(pemBytes, keyPassword)
	if err != nil {
		logrus.WithError(err).Error("parse ssh key from bytes failed")
		return nil
	}
	return ssh.PublicKeys(signer)
}

func SshRemoteRunCommand(sshClient *ssh.Client, command string) (string, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	logString := buf.String()
	//logrus.WithField("CMD:", command).Info(logString)
	if err != nil {
		return logString, fmt.Errorf("CMD: %s  OUT: %s  ERROR: %s", command, logString, err)
	}
	return logString, err
}
func NewSshClientConfig(sshUser, sshPassword, sshType, sshKey, sshKeyPassword string) (config *ssh.ClientConfig, err error) {
	if sshUser == "" {
		return nil, errors.New("ssh_user can not be empty")
	}
	config = &ssh.ClientConfig{
		Timeout:         time.Second * 3,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	switch sshType {
	case "password":
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	case "key":
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc([]byte(sshKey), []byte(sshKeyPassword))}
	default:
		return nil, fmt.Errorf("unknow ssh auth type: %s", sshType)
	}
	return
}
func CreateSimpleSshClient(sshUser, sshPassword, sshAddr string) (*ssh.Client, error) {
	if !strings.Contains(sshAddr, ":") {
		sshAddr = fmt.Sprintf("%s:22", sshAddr)
	}

	targetConfig, err := NewSshClientConfig(sshUser, sshPassword, "password", "", "")
	if err != nil {
		return nil, fmt.Errorf("cluster jumper proxy ssh config failed:%s", err)
	}
	return ssh.Dial("tcp", sshAddr, targetConfig)
}
