package model

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sshfortress/util"
)

type HardwareInfo struct {
	Disk    string `json:"hi_disk"`
	Mem     string `json:"hi_mem"`
	NetCard string `json:"hi_net_card"`
	Cpu     string `json:"hi_cpu"`
	System  string `json:"hi_system"`
	Login   string `json:"hi_login"`
	Ps      string `json:"hi_ps"`
	Port    string `json:"hi_port"`
}

func (o HardwareInfo) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *HardwareInfo) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}

func CreateHardwareInfo(machineId uint) (hi *HardwareInfo, err error) {
	hi = &HardwareInfo{}
	m := &Machine{}
	err = db.Preload("ClusterJumper").Preload("ClusterSsh").First(m, machineId).Error
	if err != nil {
		return
	}
	cs := m.ClusterSsh
	cj := m.ClusterJumper

	sshClient, err := CreateSshClientAsAdmin(m, cj, cs)
	if err != nil {
		return
	}
	defer sshClient.Close()
	hi.Disk, err = util.SshRemoteRunCommand(sshClient, "df -h")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取硬盘失败")
	}
	hi.Mem, err = util.SshRemoteRunCommand(sshClient, "free -m")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取内存失败")
	}
	hi.NetCard, err = util.SshRemoteRunCommand(sshClient, "ifconfig")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取网卡失败")
	}
	hi.Cpu, err = util.SshRemoteRunCommand(sshClient, "cat /proc/cpuinfo")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取CPU失败")
	}

	hi.System, err = util.SshRemoteRunCommand(sshClient, "uname -a;who -a;")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取系统失败")
	}
	hi.Login, err = util.SshRemoteRunCommand(sshClient, "w;last")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取Login失败")
	}
	hi.Ps, err = util.SshRemoteRunCommand(sshClient, "ps -aux")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取ps失败")
	}
	hi.Port, err = util.SshRemoteRunCommand(sshClient, "netstat -lntp")
	if err != nil {
		logrus.WithError(err).WithField("machineId", machineId).Error("获取netstat失败")
	}
	return hi, nil
}
