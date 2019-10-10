package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func SftpRename(c *gin.Context) {
	lClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	defer lClient.Close()

	sftpClient := lClient.SftpClient
	oPath := c.Query("opath")
	nPath := c.Query("npath")
	err = sftpClient.Rename(oPath, nPath)
	if handleError(c, err) {
		return
	}
	err = lClient.SaveLog("rename", fmt.Sprintf("%s => %s", oPath, nPath))
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)

}
