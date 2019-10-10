package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
	"time"
)

func Login(c *gin.Context) {
	var mdl model.User
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	ip := c.ClientIP()
	data, err := mdl.Login(ip)
	if handleError(c, err) {
		return
	}

	llogM := model.SigninLog{}
	llogM.CreatedAt = time.Now()
	llogM.ClientIp = c.ClientIP()
	llogM.UserName = mdl.Name
	llogM.Email = mdl.Email
	llogM.UserId = data.User.Id
	llogM.LoginType = "password"
	llogM.UserAgent = c.GetHeader("User-Agent")
	err = llogM.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, data)
}
