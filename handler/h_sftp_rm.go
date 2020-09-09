package handler

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
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
		err = RemoveNonemptyDirectory(sftpClient, fullPath)
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

// RemoveNonemptyDirectory removes the non empty directory in sftp server.
// sftp protocol does not allows removing non empty directory.
// we need to traverse over the file tree to remove files and directories post-orderly
func RemoveNonemptyDirectory(c *sftp.Client, path string) error {
	list, err := c.ReadDir(path)
	if err != nil {
		return err
	}

	// travarsal over the tree
	for i := 0; i < len(list); i++ {
		cur := list[i]
		newPath := filepath.Join(path, list[i].Name())
		if cur.IsDir() {
			if err := RemoveNonemptyDirectory(c, newPath); err != nil {
				return err
			}
		} else {
			if err := c.Remove(newPath); err != nil {
				return err
			}
		}
	}

	// remove current directory, which now is empty
	return c.RemoveDirectory(path)
}
