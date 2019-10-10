package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func MachineAll(c *gin.Context) {
	q := model.MachineQ{}
	err := c.ShouldBindQuery(&q)
	if handleError(c, err) {
		return
	}
	mdl := q.Machine
	query := &(q.PaginationQ)

	thisU, err := mwJwtUser(c)
	list, total, err := mdl.All(query, thisU)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}

func MachineOne(c *gin.Context) {
	var mdl model.Machine
	var err error
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	mdl.Id = id

	err = mdl.One()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func MachineCreate(c *gin.Context) {
	var mdl model.Machine
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func MachineUpdate(c *gin.Context) {
	var mdl model.Machine
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.Update()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}

func MachineDelete(c *gin.Context) {
	u, err := mwJwtUser(c)
	if handleError(c, err) {
		return
	}

	var mdl model.Machine
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}

	mdl.Id = id
	err = mdl.Delete(u)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}

//MachineHardware 获取机器的物理信息
func MachineHardware(c *gin.Context) {
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	hi, err := model.CreateHardwareInfo(id)
	if handleError(c, err) {
		return
	}
	jsonData(c, hi)
}
