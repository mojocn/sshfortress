package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type ClusterJumperQ struct {
	ClusterJumper
	PaginationQ
}

func (cq *ClusterJumperQ) All() (list *[]ClusterJumper, total uint, err error) {
	m := cq.ClusterJumper
	tx := db.Model(m)
	if m.Name != "" {
		tx = tx.Where("name like ?", "%"+m.Name+"%")
	}
	if m.SshAddr != "" {
		tx = tx.Where("ssh_addr like ?", "%"+m.SshAddr+"%")
	}
	if m.SshUser != "" {
		tx = tx.Where("ssh_user like ?", "%"+m.SshUser+"%")
	}
	list = &[]ClusterJumper{}
	total, err = crudAll(&cq.PaginationQ, tx, list)
	return
}

//ClusterJumper 私有云服务器跳板,提供ssh代理ssh服务
type ClusterJumper struct {
	BaseModel
	Name                string `json:"name" form:"name" comment:"私有云集群跳板的名称"`
	Remark              string `json:"remark" form:"remark" comment:"备注"`
	SshAddr             string `json:"ssh_addr" form:"ssh_addr"`
	SshPort             uint   `json:"ssh_port" form:"ssh_port"`
	SshUser             string `json:"ssh_user" form:"ssh_user"`
	SshType             string `json:"ssh_type" form:"ssh_type" comment:"password/key"`
	SshPassword         string `json:"-"`
	SshKey              string `json:"-" gorm:"type:text"`
	SshKeyPassword      string `json:"-"`
	InputSshPassword    string `gorm:"-" json:"input_ssh_password"`
	InputSshKeyPassword string `gorm:"-" json:"input_ssh_key_password"`
	InputSshKey         string `gorm:"-" json:"input_ssh_key"`
	Status              uint   `gorm:"default:'0'" json:"status" form:"status" comment:"机器状态 0-未知 2-连接错误 4-ssh认证错误 8-正常 "`

	Machines []Machine `json:"machines"`
}

//One
func (m *ClusterJumper) One() error {
	return crudOne(m)
}

//Update
func (m *ClusterJumper) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	if m.InputSshPassword != "" {
		m.SshPassword = m.InputSshPassword
	}
	if m.InputSshKeyPassword != "" {
		m.SshKeyPassword = m.InputSshKeyPassword
	}
	if m.InputSshKey != "" {
		m.SshKey = m.InputSshKey
	}
	return db.Model(m).Update(m).Error
}

//Create
func (m *ClusterJumper) Create() (err error) {
	m.Id = 0
	if m.InputSshPassword != "" {
		m.SshPassword = m.InputSshPassword
	}
	if m.InputSshKeyPassword != "" {
		m.SshKeyPassword = m.InputSshKeyPassword
	}
	if m.InputSshKey != "" {
		m.SshKey = m.InputSshKey
	}
	return db.Create(m).Error
}

//Delete
func (m *ClusterJumper) Delete() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	err = db.Unscoped().Where("`id` = ?", m.Id).Delete(m).Error
	return
}

//ClusterJumperBindMachine 集群跳板代理账号绑定机器json 参数绑定
type ClusterJumperBindMachine struct {
	ClusterJumperId uint   `json:"cluster_jumper_id"`
	MachineIds      []uint `json:"machine_ids"`
}

//Bind 绑定机器逻辑
func (m ClusterJumperBindMachine) Bind() error {
	if m.ClusterJumperId < 1 {
		return errors.New("cluster_ssh_id 必须大于零")
	}
	err := db.Model(Machine{}).Where("cluster_jumper_id = ?", m.ClusterJumperId).Update("cluster_jumper_id", gorm.Expr("0")).Error
	if err != nil {
		return fmt.Errorf("重置失败:%s", err)
	}
	if len(m.MachineIds) > 0 {
		return db.Model(Machine{}).Where("id in (?)", m.MachineIds).Update("cluster_jumper_id", m.ClusterJumperId).Error
	}
	return nil
}
