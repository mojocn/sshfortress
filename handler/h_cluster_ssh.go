package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func ClusterSshAll(c *gin.Context) {
	q := model.ClusterSshQ{}
	err := c.ShouldBindQuery(&q)
	if handleError(c, err) {
		return
	}
	mdl := q.ClusterSsh
	query := &(q.PaginationQ)

	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}

func ClusterSshCreate(c *gin.Context) {
	var mdl model.ClusterSsh
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

func ClusterSshUpdate(c *gin.Context) {
	var mdl model.ClusterSsh
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
func ClusterSshOne(c *gin.Context) {
	var mdl model.ClusterSsh
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

func ClusterSshBindMachines(c *gin.Context) {
	var mdl model.ClusterSshBindMachine
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

func ClusterSshDelete(c *gin.Context) {
	var mdl model.ClusterSsh
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
