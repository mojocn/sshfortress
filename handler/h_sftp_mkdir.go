package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
	"strings"
)

func SftpMkdir(c *gin.Context) {

	lClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	sftpClient := lClient.SftpClient
	fullPath := c.Query("path")
	if strings.HasPrefix(fullPath, "$HOME") {
		wd, err := sftpClient.Getwd()
		if handleError(c, err) {
			return
		}
		fullPath = strings.Replace(fullPath, "$HOME", wd, 1)
	}
	err = sftpClient.Mkdir(fullPath)
	if handleError(c, err) {
		return
	}
	err = lClient.SaveLog("mkdir", fullPath)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}

func getSftpClient(c *gin.Context) (sc *model.LogicSftpClient, err error) {
	user, err := mwJwtUser(c)
	if err != nil {
		return nil, err
	}
	idx, err := parseParamID(c)
	if err != nil {
		return nil, err
	}
	return model.CreateLogicSftpClient(user, idx, c.ClientIP())
}
