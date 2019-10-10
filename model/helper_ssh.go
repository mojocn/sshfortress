package model

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"sshfortress/util"
)

func CreateSshClientAsAdmin(m *Machine, cj *ClusterJumper, cs *ClusterSsh) (c *ssh.Client, err error) {
	if m.Id < 1 {
		return nil, errors.New("CreateSshClientAsAdmin m is not valid")
	}
	targetConfig, err := util.NewSshClientConfig(cs.SshUser, cs.SshPassword, cs.SshType, cs.SshKey, cs.SshKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("cluster jumper proxy ssh config failed:%s", err)
	}
	targetAddr := fmt.Sprintf("%s:%d", m.SshIp, m.SshPort)

	var proxyConfig *ssh.ClientConfig
	var proxyAddr string
	if cj != nil && cj.Id > 0 {
		//使用私有云集群跳板登陆
		proxyConfig, err = util.NewSshClientConfig(cj.SshUser, cj.SshPassword, cj.SshType, cj.SshKey, cj.SshKeyPassword)
		if err != nil {
			return nil, fmt.Errorf("cluster jumper proxy ssh config failed:%s", err)
		}
		proxyAddr = fmt.Sprintf("%s:%d", cj.SshAddr, cj.SshPort)
	}
	c, err = createSshProxySshClient(targetConfig, proxyConfig, targetAddr, proxyAddr)
	return
}

func CreateSshClientAsUser(m *Machine, user *User, cj *ClusterJumper) (c *ssh.Client, err error) {
	targetConfig, err := util.NewSshClientConfig(user.SshUserName(), user.SshPassword, "password", "", "")
	if err != nil {
		return nil, fmt.Errorf("cluster jumper proxy ssh config failed:%s", err)
	}
	targetAddr := fmt.Sprintf("%s:%d", m.SshIp, m.SshPort)

	var proxyConfig *ssh.ClientConfig
	var proxyAddr string
	if cj != nil && cj.Id > 0 {
		//使用私有云集群跳板登陆
		proxyConfig, err = util.NewSshClientConfig(cj.SshUser, cj.SshPassword, cj.SshType, cj.SshKey, cj.SshKeyPassword)
		if err != nil {
			return nil, fmt.Errorf("cluster jumper proxy ssh config failed:%s", err)
		}
		proxyAddr = fmt.Sprintf("%s:%d", cj.SshAddr, cj.SshPort)
	}
	c, err = createSshProxySshClient(targetConfig, proxyConfig, targetAddr, proxyAddr)
	return
}

func createSshProxySshClient(targetSshConfig, proxySshConfig *ssh.ClientConfig, targetAddr, proxyAddr string) (client *ssh.Client, err error) {
	if proxySshConfig == nil {
		return ssh.Dial("tcp", targetAddr, targetSshConfig)
	}

	proxyClient, err := ssh.Dial("tcp", proxyAddr, proxySshConfig)
	if err != nil {
		return
	}
	conn, err := proxyClient.Dial("tcp", targetAddr)
	if err != nil {
		return
	}
	ncc, chans, reqs, err := ssh.NewClientConn(conn, targetAddr, targetSshConfig)
	if err != nil {
		return
	}
	client = ssh.NewClient(ncc, chans, reqs)
	return
}
