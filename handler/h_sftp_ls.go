package handler

import (
	"github.com/gin-gonic/gin"
	"path"
	"time"
)

type Ls struct {
	Name  string    `json:"name"`
	Path  string    `json:"path"` // including Name
	Size  int64     `json:"size"`
	Time  time.Time `json:"time"`
	Mod   string    `json:"mod"`
	IsDir bool      `json:"is_dir"`
}

func SftpLs(c *gin.Context) {
	lClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	defer lClient.Close()

	sftpClient := lClient.SftpClient
	dirPath := c.DefaultQuery("path", "$HOME")
	if dirPath == "$HOME" {
		wd, err := sftpClient.Getwd()
		if handleError(c, err) {
			return
		}
		dirPath = wd
	}
	files, err := sftpClient.ReadDir(dirPath)
	if handleError(c, err) {
		return
	}
	fileList := make([]Ls, 0) // this will not be converted to null if slice is empty.
	for _, file := range files {
		tt := Ls{Name: file.Name(), Size: file.Size(), Path: path.Join(dirPath, file.Name()), Time: file.ModTime(), Mod: file.Mode().String(), IsDir: file.IsDir()}
		fileList = append(fileList, tt)
	}
	//记录sftp 日志
	err = lClient.SaveLog("ls", dirPath)
	if handleError(c, err) {
		return
	}
	jsonData(c, fileList)
}
