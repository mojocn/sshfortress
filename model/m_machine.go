package model

import (
	"errors"
)

type MachineQ struct {
	Machine
	PaginationQ
}

type Machine struct {
	BaseModel
	Name            string         `gorm:"type:varchar(50);unique_index" json:"name" form:"name"`
	SshIp           string         `json:"ssh_ip" form:"ssh_ip"`
	SshPort         uint           `json:"ssh_port"`
	LanIp           string         `json:"lan_ip" form:"lan_ip"`
	WanIp           string         `json:"wan_ip" form:"wan_ip"`
	Cate            uint           `gorm:"default:'2'" json:"cate" comment:"机器性质:2:无外网ip 4:外网可以访问"`
	ClusterSshId    uint           `gorm:"index" json:"cluster_ssh_id" comment:"关联集群管理ssh账号"`
	ClusterJumperId *uint          `gorm:"index,default:'0'" json:"cluster_jumper_id" form:"cluster_jumper_id" comment:"集群代理ID 关联clusterJumper表"`
	UserId          uint           `gorm:"index" json:"user_id" comment:"机器的添加者"`
	Status          uint           `gorm:"default:'0'" json:"status" form:"status" comment:"机器状态 0-未知 2-连接错误 4-ssh认证错误 8-正常 "`
	User            User           `gorm:"association_autoupdate:false;association_autocreate:false" json:"user"`
	ClusterSsh      *ClusterSsh    `gorm:"association_autoupdate:false;association_autocreate:false" json:"cluster_ssh,omitempty"`
	ClusterJumper   *ClusterJumper `gorm:"association_autoupdate:false;association_autocreate:false" json:"cluster_jumper,omitempty"`
	//Hardware        HardwareInfo   `gorm:"type:json" json:"hardware"`
}

func (m *Machine) AfterFind() (err error) {
	return
}

//One
func (m *Machine) One() error {
	return crudOne(m)
}

//All
func (m Machine) All(q *PaginationQ, user *User) (list *[]Machine, total uint, err error) {
	tx := db.Model(m).Preload("ClusterSsh").Preload("ClusterJumper") //.Where("ancestor_path like ?", m.qAncetorPath())
	list = &[]Machine{}
	//role ==2
	//显示全部的机器
	if m.Name != "" {
		tx = tx.Where("`name` like ?", "%"+m.Name+"%")
	}
	if m.SshIp != "" {
		tx = tx.Where("`ssh_ip` like ?", "%"+m.SshIp+"%")
	}
	if m.WanIp != "" {
		tx = tx.Where("`wan_ip` like ?", "%"+m.WanIp+"%")
	}
	if m.LanIp != "" {
		tx = tx.Where("`lan_ip` like ?", "%"+m.LanIp+"%")
	}

	if user.Role == 4 {
		//普通用户显示自由机器和授权的机器
		machineIds := []uint{}
		err = db.Model(MachineUser{}).Where("user_id = ?", user.Id).Pluck("machine_id", &machineIds).Error
		if err != nil {
			return nil, 0, err
		}
		if len(machineIds) > 0 {
			tx = tx.Where("`user_id` = ? OR `id` in (?)", user.Id, machineIds)
		} else {
			tx = tx.Where("`user_id` = ?", user.Id)
		}
	}

	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *Machine) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}

	return db.Model(m).Update(m).Error
}

//Create
func (m *Machine) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

//Delete
func (m *Machine) Delete(u *User) (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	//删除用户与机器的关联
	err = db.Where("machine_id = ?", m.Id).Delete(MachineUser{}).Error
	if err != nil {
		return
	}
	if u.Role == UserRoleAdmin {
		return crudDelete(m)
	}
	err = db.Unscoped().Where("`id` = ? AND `user_id` = ?", m.Id, u.Id).Delete(m).Error
	return
}
