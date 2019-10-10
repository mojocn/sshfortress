package model

import (
	"errors"
)

type LogicMachineUser struct {
	UserId    uint          `json:"user_id"`
	MachineId uint          `json:"machine_id"`
	Relations []MachineUser `json:"relations"`
}

//UserBindMachines 用户绑定机器
func (l LogicMachineUser) UserBindMachines() (err error) {
	if l.UserId < 1 {
		err = errors.New("user_id 必须未非零值")
		return
	}
	banRelations := []MachineUser{}
	err = db.Where("user_id =?", l.UserId).Find(&banRelations).Error
	if err != nil {
		//关系表中机器已经被删除, 脏数据库
		return
	}

	for _, v := range banRelations {
		err = v.DeleteMachineUserRelation()
		if err != nil {
			return
		}
	}

	for _, v := range l.Relations {
		v.UserId = l.UserId
		err = v.CreateSshUserOnMachineSaveRelation()
		if err != nil {
			return
		}
	}
	return
}

//UserBindMachines 机器绑定用户
func (l LogicMachineUser) MachineBindUsers() (err error) {
	if l.MachineId < 1 {
		err = errors.New("user_id 必须未非零值")
		return
	}
	banRelations := []MachineUser{}
	err = db.Where("machine_id =?", l.MachineId).Find(&banRelations).Error
	if err != nil {
		return
	}

	for _, v := range banRelations {
		//暂时禁用ssh账号在目标机器上登陆
		err = v.DeleteMachineUserRelation()
		if err != nil {
			return
		}
	}

	for _, v := range l.Relations {
		v.MachineId = l.MachineId
		err = v.CreateSshUserOnMachineSaveRelation()
		if err != nil {
			return
		}
	}
	return
}
