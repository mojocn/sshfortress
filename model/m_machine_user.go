package model

import "github.com/jinzhu/gorm"

//MachineUser 机器和用户关系表
type MachineUser struct {
	Id        uint `gorm:"primary_key" json:"id"`
	MachineId uint `gorm:"index" json:"machine_id"`
	UserId    uint `gorm:"index" json:"user_id"`
	SudoType  uint `gorm:"default:2" json:"sudo_type" comment:"ssh账号权限  2-普通 4-sudo密码 8-sudo免密码"`
}

func (r MachineUser) GetMachineIdsBy(userId uint) (list []MachineUser, err error) {
	err = db.Model(r).Where("user_id = ?", userId).Find(&list).Error
	return
}
func (r MachineUser) GetUserIdsBy(machineId uint) (list []MachineUser, err error) {
	err = db.Model(r).Where("machine_id = ?", machineId).Find(&list).Error
	return
}

//CreateSshUserOnMachineSaveRelation 在目标机器上创建ssh账号同时成功就写入关系到数据库中
func (r *MachineUser) CreateSshUserOnMachineSaveRelation() error {
	l, err := CreateLogicMachineSshAccount(r)
	if err != nil {
		return err
	}
	defer l.Close()
	err = l.AuthSshUserToMachine()
	if err != nil {
		return err
	}
	return db.Create(r).Error
}

//Create 在目标机器上创建ssh账号同时成功就写入关系到数据库中
func (r *MachineUser) DeleteMachineUserRelation() error {
	l, err := CreateLogicMachineSshAccount(r)
	if err == gorm.ErrRecordNotFound {
		return crudDelete(r)
	}
	if err != nil {
		return err
	}
	defer l.Close()
	err = l.LockSshAccount()
	if err != nil {
		return err
	}
	return crudDelete(r)
}
