package model

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

//CreateLogicMachineSshAccount 创建LogicMachineSshAccount
func CreateLogicMachineSshAccount(r *MachineUser) (sa *LogicMachineSshAccount, err error) {
	user := &User{}
	machine := &Machine{}
	cs := &ClusterSsh{}
	cj := &ClusterJumper{}
	//检验userID 是否有效
	err = db.First(user, r.UserId).Error
	if err != nil {
		return
	}
	err = db.First(machine, r.MachineId).Related(cs, "ClusterSsh").Related(cj, "ClusterJumper").Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	sa = &LogicMachineSshAccount{Machine: machine, User: user, ClusterSsh: cs, MachineUser: r}

	sa.sshClient, err = CreateSshClientAsAdmin(machine, cj, cs)
	if err != nil {
		return nil, fmt.Errorf("create cluster jumper proxy ssh client for ssh account failed:%s", err)
	}

	if sa.User.SshPassword == "" {
		return nil, fmt.Errorf("用户%s sshPassword 不能为空", sa.User.Email)
	}
	sa.sshUser = sa.User.SshUserName()
	sa.sshPassword = sa.User.SshPassword
	sa.sudoType = sa.MachineUser.SudoType

	return
}

//LogicMachineSshAccount 处理目标机器中ssh账号逻辑,设计到ssh 远程代码执行
type LogicMachineSshAccount struct {
	sshUser     string
	sshPassword string
	sudoType    uint
	MachineUser *MachineUser
	Machine     *Machine
	ClusterSsh  *ClusterSsh
	User        *User
	sshClient   *ssh.Client
}

//Close 关闭sshClient 和日志处理
func (msa *LogicMachineSshAccount) Close() {
	if msa.sshClient != nil {
		err := msa.sshClient.Close()
		if err != nil {
			logrus.WithError(err).Error("LogicMachineSshAccount.Close 关闭sshClient 失败")
		}
	}
}

//AuthSshUserToMachine 用户与机器授权的调用这个方法在目标机器上创建/禁用账号
func (msa *LogicMachineSshAccount) AuthSshUserToMachine() (err error) {
	err = msa.CheckSshAccountExist()
	if err != nil {
		logrus.WithError(err).WithField("sshUser", msa.sshUser).Info("用户名不存在,then 创建用户")
		err = msa.CreateSshUser()
		if err != nil {
			return
		}
	}
	err = msa.UnlockSshAccount()
	if err != nil {
		return
	}
	err = msa.ChangePassword()
	if err != nil {
		return
	}
	err = msa.SetSudoType()
	if err != nil {
		return
	}
	return
}

func (msa *LogicMachineSshAccount) runCmd(command string) (string, error) {
	session, err := msa.sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	logString := buf.String()
	logrus.WithField("CMD:", command).Info(logString)
	if err != nil {
		return logString, fmt.Errorf("命令: %s  输出: %s  错误: %s,优先检查机器绑定的集群管理SSH账号是否有sudo NOPASSWD权限", command, logString, err)
	}
	return logString, nil
}

//CheckSshAccountExist 检查目标机器上ssh用户是否存在
func (msa *LogicMachineSshAccount) CheckSshAccountExist() error {
	cmd := fmt.Sprintf("sudo id -u '%s'", msa.sshUser)
	_, err := msa.runCmd(cmd)
	return err
}

//LockSshAccount 锁定目标ssh账号禁止登陆,相当于软删除
func (msa *LogicMachineSshAccount) LockSshAccount() error {
	cmd := fmt.Sprintf("sudo usermod -L '%s' || true", msa.sshUser)
	_, err := msa.runCmd(cmd)
	return err
}

//UnlockSshAccount 解禁目标ssh账号禁止登陆,相当于恢复
func (msa *LogicMachineSshAccount) UnlockSshAccount() error {
	cmd := fmt.Sprintf("sudo usermod -U '%s' || true", msa.sshUser)
	_, err := msa.runCmd(cmd)
	return err
}

//ChangePassword  修改ssh账号密码
func (msa *LogicMachineSshAccount) ChangePassword() error {
	cmd := fmt.Sprintf("echo '%s:%s' | sudo chpasswd", msa.sshUser, msa.sshPassword)
	_, err := msa.runCmd(cmd)
	return err
}

//RemoveSudo 移除sudo权限
func (msa *LogicMachineSshAccount) RemoveSudo() error {
	//清楚 /etc/sudoers 文件中 用户名信息
	cmd := fmt.Sprintf(`sudo sed -i '/^%s\b/d' /etc/sudoers`, msa.sshUser)
	_, err := msa.runCmd(cmd)
	return err
}

//RemoveSshAccountHard 删除用户的目相关文件
func (msa *LogicMachineSshAccount) RemoveSshAccountHard() error {
	//清楚 /etc/sudoers 文件中 用户名信息
	cmd := fmt.Sprintf(`sudo userdel -r -f -Z '%s'`, msa.sshUser)
	_, err := msa.runCmd(cmd)
	return err
}

//SetSudoType 设置sudo 同时是否免书密码
func (msa *LogicMachineSshAccount) SetSudoType() error {
	err := msa.RemoveSudo()
	if err != nil {
		return err
	}
	switch msa.sudoType {
	case 4:
		return msa.AddSudo(false)
	case 8:
		return msa.AddSudo(true)
	default:
		return nil
	}
}

//RemoveSudo 移除sudo权限 isNeedPassword sudo是否需要密码
func (msa *LogicMachineSshAccount) AddSudo(isNoPassword bool) error {
	var cmd string
	//https://stackoverflow.com/questions/20895619/why-cant-i-echo-contents-into-a-new-file-as-sudo
	if isNoPassword {
		cmd = fmt.Sprintf(`sudo bash -c "echo '%s   ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers;"`, msa.sshUser)
	} else {
		cmd = fmt.Sprintf(`sudo bash -c "echo '%s   ALL=(ALL) ALL' >> /etc/sudoers;"`, msa.sshUser)
	}
	_, err := msa.runCmd(cmd)
	return err
}

//CreateSshUser 创建ssh账号
func (msa *LogicMachineSshAccount) CreateSshUser() error {
	cmd := fmt.Sprintf("sudo useradd -m '%s'", msa.sshUser)
	_, err := msa.runCmd(cmd)
	return err
}
