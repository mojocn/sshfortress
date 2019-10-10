package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"time"
)

type LogicSshClient struct {
	SshClient *ssh.Client
	User      *User
	Machine   *Machine
	SshUser   string `comment:"有可能从用户表中读取, 也有可能冲账号表中读取"`
	StartTime time.Time
	ClientIp  string
}

//CreateLogicSshClient 创建ssh client
func CreateLogicSshClient(user *User, machineId uint, clientIp string) (client *LogicSshClient, err error) {
	client = &LogicSshClient{}
	client.ClientIp = clientIp
	client.StartTime = time.Now()
	client.User = user
	m := &Machine{}
	if user.IsAdmin() {
		//管理用户
		//使用集群账号登陆,不检查机器授权关系
		err = db.Preload("ClusterSsh").Preload("ClusterJumper").First(m, machineId).Error
		if err != nil {
			return
		}
		cj := m.ClusterJumper
		cs := m.ClusterSsh
		client.SshClient, err = CreateSshClientAsAdmin(m, cj, cs)
		if err != nil {
			return nil, fmt.Errorf("create cluster jumper proxy ssh client for admin failed:%s", err)
		}
		client.SshUser = cs.SshUser
	} else {
		//普通用户
		//使用自身账号登陆,检查机器授权关系
		relation := MachineUser{}
		if db.Where("user_id = ?", user.Id).Where("machine_id = ?", machineId).First(&relation).RecordNotFound() {
			err = fmt.Errorf("用户Id: %d 没有ssh访问机器Id: %d 权限", user.Id, machineId)
			return
		}
		err = db.Preload("ClusterJumper").First(m, machineId).Error
		if err != nil {
			return
		}

		client.SshClient, err = CreateSshClientAsUser(m, user, m.ClusterJumper)
		if err != nil {
			return nil, fmt.Errorf("create cluster jumper proxy ssh client for user failed:%s", err)
		}
		client.SshUser = user.SshUserName()
	}
	client.Machine = m
	return
}

//SaveLog 保存ssh连接日志
func (l *LogicSshClient) SaveLog(isFlagged, hasEditor bool, tl JsonArrayString) (err error) {
	xtermLog := SshLog{
		SshUser:   l.SshUser,
		StartedAt: l.StartTime,
		UserId:    l.User.Id,
		MachineId: l.Machine.Id,
		Status:    0,
		ClientIp:  l.ClientIp,
		TLog:      tl,
	}
	if isFlagged {
		xtermLog.Status = 8
	}
	if hasEditor {
		//session中进行了文本编辑
		xtermLog.Status = 32
	}
	return xtermLog.Create()
}

//Close 关闭连接
func (l *LogicSshClient) Close() {
	if l.SshClient == nil {
		return
	}
	err := l.SshClient.Close()
	if err != nil {
		logrus.WithError(err).Error("close LogicSshClient.SshClient failed")
	}
}
