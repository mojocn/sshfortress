package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func ConfigAll(c *gin.Context) {
	q := model.ConfigQ{}
	err := c.ShouldBindQuery(&q)
	if handleError(c, err) {
		return
	}
	mdl := q.Config
	query := &(q.PaginationQ)

	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}

func ConfigUpdate(c *gin.Context) {
	var mdl model.Config
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
