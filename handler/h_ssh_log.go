package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func SshLogAll(c *gin.Context) {
	u, err := mwJwtUser(c)
	if handleError(c, err) {
		return
	}
	query := &model.SshLogQ{}
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

func SshLogUpdate(c *gin.Context) {
	var mdl model.SshLog
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

func SshLogDelete(c *gin.Context) {
	ids := []int{}
	err := c.ShouldBind(&ids)
	if handleError(c, err) {
		return
	}
	var mdl model.SshLog
	err = mdl.Delete(ids)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
