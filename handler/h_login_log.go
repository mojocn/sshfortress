package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func SigninLogAll(c *gin.Context) {
	q := model.SigninLogQ{}
	err := c.ShouldBindQuery(&q)
	if handleError(c, err) {
		return
	}
	query := &(q.PaginationQ)
	list, total, err := q.Search(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}

func SigninLogDelete(c *gin.Context) {
	ids := []int{}
	err := c.ShouldBind(&ids)
	if handleError(c, err) {
		return
	}
	var mdl model.SigninLog
	err = mdl.Delete(ids)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
