package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func SshFilterGroupAll(c *gin.Context) {
	var ugq model.SshFilterGroupQ
	err := c.ShouldBindQuery(&ugq)
	if handleError(c, err) {
		return
	}
	page, err := ugq.Search()
	if handleError(c, err) {
		return
	}
	c.JSON(200, page)
}

func SshFilterGroupCreate(c *gin.Context) {
	var mdl model.SshFilterGroup
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

func SshFilterGroupUpdate(c *gin.Context) {
	var mdl model.SshFilterGroup
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
func SshFilterGroupOne(c *gin.Context) {
	var mdl model.SshFilterGroup
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

func SshFilterGroupDelete(c *gin.Context) {
	var mdl model.SshFilterGroup
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
