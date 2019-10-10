package model

import (
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"time"
)

type LogicSftpClient struct {
	User       *User
	Machine    *Machine
	SshUser    string `comment:"有可能从用户表中读取, 也有可能冲账号表中读取"`
	SftpClient *sftp.Client
	ClientIp   string
}

func CreateLogicSftpClient(user *User, machineId uint, ip string) (client *LogicSftpClient, err error) {
	client = &LogicSftpClient{}
	client.ClientIp = ip
	client.User = user
	scl, err := CreateLogicSshClient(user, machineId, ip)
	if err != nil {
		return nil, err
	}
	client.Machine = scl.Machine
	client.SshUser = scl.SshUser
	client.SftpClient, err = sftp.NewClient(scl.SshClient, sftp.MaxPacket(maxPacket))
	if err != nil {
		return nil, err
	}
	return
}

const maxPacket = 1024 * 32

func (l *LogicSftpClient) SaveLog(action, path string) (err error) {
	m := SftpLog{}
	m.SshUser = l.SshUser
	m.UserId = l.User.Id
	m.MachineId = l.Machine.Id
	m.ClientIp = l.ClientIp
	m.Action = action
	m.Path = path
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	switch action {
	case "rm-dir", "rm-file":
		m.Status = 16
	case "upload":
		m.Status = 8
	case "rename":
		m.Status = 4
	}
	return m.Create()
}

func (l *LogicSftpClient) Close() {
	if l.SftpClient == nil {
		return
	}
	err := l.SftpClient.Close()
	if err != nil {
		logrus.WithError(err).Error("close LogicSftpClient.SftpClient failed")
	}
	return
}
