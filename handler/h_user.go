package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func UserAll(c *gin.Context) {
	q := model.UserQ{}
	err := c.ShouldBindQuery(&q)
	if handleError(c, err) {
		return
	}
	query := &(q.PaginationQ)
	mdl := q.User
	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}

func UserCreate(c *gin.Context) {
	var mdl model.User
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	u, err := mwJwtUser(c)
	if handleError(c, err) {
		return
	}
	mdl.ParentId = u.Id
	mdl.AncestorPath = fmt.Sprintf("%s/%d", u.AncestorPath, u.Id)
	err = mdl.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func UserUpdate(c *gin.Context) {
	user, err := mwJwtUser(c)
	if handleError(c, err) {
		return
	}
	var mdl model.User
	err = c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	if !user.IsAdmin() && mdl.Id != user.Id {
		jsonError(c, "common user only can update himself/herself")
		return
	}

	err = mdl.Update()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}
func UserOne(c *gin.Context) {
	var mdl model.User
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
func UserDelete(c *gin.Context) {
	var mdl model.User
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	uid, err := mwJwtUid(c)
	if handleError(c, err) {
		return
	}
	if uid == id {
		jsonError(c, "can't delete your own account")
		return
	}
	mdl.Id = id
	err = mdl.Delete()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
