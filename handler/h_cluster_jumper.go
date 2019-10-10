package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func ClusterJumperAll(c *gin.Context) {
	q := model.ClusterJumperQ{}
	err := c.ShouldBindQuery(&q)
	if handleError(c, err) {
		return
	}
	list, total, err := q.All()
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, &q.PaginationQ)
}

func ClusterJumperCreate(c *gin.Context) {
	var mdl model.ClusterJumper
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

func ClusterJumperUpdate(c *gin.Context) {
	var mdl model.ClusterJumper
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
func ClusterJumperOne(c *gin.Context) {
	var mdl model.ClusterJumper
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

func ClusterJumperDelete(c *gin.Context) {
	var mdl model.ClusterJumper
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	mdl.Id = id
	err = mdl.Delete()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}

func ClusterJumperBindMachines(c *gin.Context) {
	var mdl model.ClusterJumperBindMachine
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.Bind()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
