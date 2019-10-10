package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SftpRm(c *gin.Context) {
	lClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	defer lClient.Close()
	sftpClient := lClient.SftpClient
	fullPath := c.Query("path")
	dirOrFile := c.Query("dirOrFile")

	if fullPath == "/" || fullPath == "$HOME" {
		jsonError(c, "can't delete / or $HOME dir")
		return
	}
	logAction := ""
	switch dirOrFile {
	case "dir":
		//sftp 删除非空文件夹错误
		//todo:: 解决方案使用rm -rf 命令来删除
		err = sftpClient.RemoveDirectory(fullPath)
		if err != nil {
			jsonError(c, fmt.Sprintf("can't delete none empty directory: %s", err))
			return
		}
		logAction = "rm-dir"
	case "file":
		err = sftpClient.Remove(fullPath)
		logAction = "rm-file"
	default:
		err = errors.New("dirOrFile 参数是必须的且是file/dir")
	}
	if handleError(c, err) {
		return
	}
	err = lClient.SaveLog(logAction, fullPath)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)

}
