package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type ClusterSshQ struct {
	PaginationQ
	ClusterSsh
}

//ClusterSsh
type ClusterSsh struct {
	BaseModel
	SshUser             string    `json:"ssh_user" form:"ssh_user"`
	SshType             string    `json:"ssh_type" form:"ssh_type" comment:"password/key"`
	Remark              string    `json:"remark" form:"remark"`
	SshPassword         string    `json:"-"`
	SshKey              string    `json:"-" gorm:"type:text"`
	SshKeyPassword      string    `json:"-"`
	InputSshPassword    string    `gorm:"-" json:"input_ssh_password"`
	InputSshKeyPassword string    `gorm:"-" json:"input_ssh_key_password"`
	InputSshKey         string    `gorm:"-" json:"input_ssh_key"`
	Machines            []Machine `json:"machines"`
}

//ClusterSshBindMachine 集群账号绑定机器json 参数绑定
type ClusterSshBindMachine struct {
	ClusterSshId uint   `json:"cluster_ssh_id"`
	MachineIds   []uint `json:"machine_ids"`
}

//Bind 绑定机器逻辑
func (m ClusterSshBindMachine) Bind() error {
	if m.ClusterSshId < 1 {
		return errors.New("cluster_ssh_id 必须大于零")
	}
	err := db.Model(Machine{}).Where("cluster_ssh_id = ?", m.ClusterSshId).Update("cluster_ssh_id", gorm.Expr("NULL")).Error
	if err != nil {
		return fmt.Errorf("重置失败:%s", err)
	}
	if len(m.MachineIds) > 0 {
		return db.Model(Machine{}).Where("id in (?)", m.MachineIds).Update("cluster_ssh_id", m.ClusterSshId).Error
	}
	return nil
}

func (m *ClusterSsh) AfterFind() (err error) {
	return
}

//One
func (m *ClusterSsh) One() error {
	return db.Preload("Machines").First(m).Error
}

//All
func (m ClusterSsh) All(q *PaginationQ) (list *[]ClusterSsh, total uint, err error) {
	list = &[]ClusterSsh{}
	tx := db.Model(m) //.Where("ancestor_path like ?", m.qAncetorPath())
	if m.Remark != "" {
		tx = tx.Where("remark like ?", "%"+m.Remark+"%")
	}
	if m.SshUser != "" {
		tx = tx.Where("ssh_user like ?", "%"+m.SshUser+"%")
	}
	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *ClusterSsh) Update() (err error) {
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
func (m *ClusterSsh) Create() (err error) {
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
func (m *ClusterSsh) Delete() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	return crudDelete(m)
}
