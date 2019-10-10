package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

//MachineUserMachineIds 根据user_id 获取 machine_ids array
func MachineUserMachineIds(c *gin.Context) {
	userId, err := queryUint(c, "user_id")
	if handleError(c, err) {
		return
	}
	list, err := (model.MachineUser{}).GetMachineIdsBy(userId)
	if handleError(c, err) {
		return
	}
	jsonData(c, list)
}

//MachineUserUserIds machine_id 获取 user_ids array
func MachineUserUserIds(c *gin.Context) {
	machineId, err := queryUint(c, "machine_id")
	if handleError(c, err) {
		return
	}
	list, err := (model.MachineUser{}).GetUserIdsBy(machineId)
	if handleError(c, err) {
		return
	}
	jsonData(c, list)
}

//MachineUserBindUsers 选择一台机器给他绑定用户
func MachineUserBindUsers(c *gin.Context) {
	var q model.LogicMachineUser
	err := c.ShouldBind(&q)
	if handleError(c, err) {
		return
	}
	err = q.MachineBindUsers()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}

//MachineUserBindMachines 选择一个用户给他绑定机器
func MachineUserBindMachines(c *gin.Context) {
	var q model.LogicMachineUser
	err := c.ShouldBind(&q)
	if handleError(c, err) {
		return
	}
	err = q.UserBindMachines()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
