package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sshfortress/model"
	"sshfortress/stat"
	"strconv"
	"time"
)

func jsonError(c *gin.Context, msg interface{}) {
	stat.GaugeVecApiError.WithLabelValues("API").Inc()
	var ms string
	switch v := msg.(type) {
	case string:
		ms = v
	case error:
		ms = v.Error()
	default:
		ms = ""
	}
	c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": ms})
}
func jsonAuthError(c *gin.Context, msg interface{}) {
	c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"ok": false, "msg": msg})
}

func jsonData(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": data})
}

//func jsonPagination(c *gin.Context, list interface{}, total uint, query *model.PaginationQ) {
//	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": list, "total": total, "offset": query.Offset, "limit": query.Size})
//}
func jsonSuccess(c *gin.Context) {
	c.AbortWithStatusJSON(200, gin.H{"ok": true, "msg": "success"})
}
func jsonPagination(c *gin.Context, list interface{}, total uint, query *model.PaginationQ) {
	c.JSON(200, gin.H{"ok": true, "data": list, "total": total, "page": query.Page, "size": query.Size})
}
func handleError(c *gin.Context, err error) bool {
	if err != nil {
		//logrus.WithError(err).Error("gin context http handler error")
		jsonError(c, err.Error())
		return true
	}
	return false
}
func handlerAuthMiddlewareError(c *gin.Context, err error) bool {
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"msg": err.Error()})
		return true
	}
	return false
}
func wshandleError(err error) bool {
	if err != nil {
		stat.GaugeVecApiError.WithLabelValues("WS").Inc()
		logrus.WithError(err).Error("handler ws ERROR:")
		return true
	}
	return false
}
func wshandleErrorPro(wc *websocket.Conn, err error) bool {
	if err != nil {
		logrus.WithError(err).Error("ssh-websocket error")
		wc.WriteControl(websocket.CloseMessage, []byte(err.Error()), time.Now())
		return true
	}
	return false
}

func parseParamID(c *gin.Context) (uint, error) {
	id := c.Param("id")
	parseId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, errors.New("id must be an unsigned int")
	}
	return uint(parseId), nil
}

func mwJwtUser(c *gin.Context) (*model.User, error) {
	uid, err := mwJwtUid(c)
	if err != nil {
		return nil, err
	}
	sessionUserKey := fmt.Sprintf("thisUser:%d", uid)
	vu, ok := c.Get(sessionUserKey)
	if ok {
		tu, okay := vu.(model.User)
		if okay {
			return &tu, nil
		} else {
			return nil, fmt.Errorf("gin.context %s 不是model.User", sessionUserKey)
		}
	}
	//session context 没有值则mysql数据库中取,并设置context 值
	user := model.User{}
	user.Id = uid
	err = user.One()
	if err != nil {
		return nil, fmt.Errorf("context can not get user of %d,error:%s", user.Id, err)
	}
	c.Set(sessionUserKey, user)
	return &user, nil
}
func mwJwtUid(c *gin.Context) (uint, error) {
	return getCtxUint(c, jwtCtxUidKey)
}

func getCtxUint(c *gin.Context, key string) (uint, error) {
	v, exist := c.Get(key)
	if !exist {
		return 0, fmt.Errorf("context has no value for %s", key)
	}
	uintV, ok := v.(uint)
	if ok {
		return uintV, nil
	}
	return 0, fmt.Errorf("key for %s in gin.Context value is %v but not a uint type", key, v)
}

func queryUint(c *gin.Context, key string) (v uint, err error) {
	sv, ok := c.GetQuery(key)
	if !ok {
		err = fmt.Errorf("query of %s is not exist", key)
		return
	}
	parseId, err := strconv.ParseUint(sv, 10, 32)
	if err != nil {
		err = fmt.Errorf("query of %s is not a uint", key)
		return
	}
	return uint(parseId), nil
}

func checkJwtUserAdmin(c *gin.Context) (u *model.User, err error) {
	u, err = mwJwtUser(c)
	if err != nil {
		return
	}
	if u.Role != model.UserRoleAdmin {
		return nil, errors.New("没有管理员权限")
	}
	return
}
