package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func FeedbackAll(c *gin.Context) {
	q := model.FeedbackQ{}

	err := c.ShouldBindQuery(&q)
	if handleError(c, err) {
		return
	}
	mdl := q.Feedback
	query := &(q.PaginationQ)
	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}

func FeedbackOne(c *gin.Context) {
	var mdl model.Feedback
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

func FeedbackCreate(c *gin.Context) {
	var mdl model.Feedback
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	uid, err := mwJwtUid(c)
	if handleError(c, err) {
		return
	}
	mdl.UserId = uid
	err = mdl.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func FeedbackUpdate(c *gin.Context) {
	var mdl model.Feedback
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

func FeedbackDelete(c *gin.Context) {
	ids := []int{}
	err := c.ShouldBind(&ids)
	if handleError(c, err) {
		return
	}
	var mdl model.Feedback
	err = mdl.Delete(ids)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
