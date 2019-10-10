package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func SftpLogAll(c *gin.Context) {
	u, err := mwJwtUser(c)
	if handleError(c, err) {
		return
	}
	query := &model.SftpLogQ{}
	err = c.ShouldBindQuery(query)
	if handleError(c, err) {
		return
	}
	list, total, err := query.Search(u)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, &query.PaginationQ)
}

func SftpLogUpdate(c *gin.Context) {
	var mdl model.SftpLog
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

func SftpLogDelete(c *gin.Context) {
	ids := []int{}
	err := c.ShouldBind(&ids)
	if handleError(c, err) {
		return
	}
	var mdl model.SftpLog
	err = mdl.Delete(ids)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
