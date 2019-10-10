package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"io/ioutil"
	"net/http"
	"time"
)

func SftpDl(c *gin.Context) {
	fullPath := c.Query("path")
	sourceType := c.Query("type") //file or dir
	lClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	defer lClient.Close()
	sftpClient := lClient.SftpClient
	if sourceType == "file" {
		fi, err := sftpClient.Stat(fullPath)
		if handleError(c, err) {
			return
		}
		f, err := sftpClient.Open(fullPath)
		defer f.Close()
		if handleError(c, err) {
			return
		}
		//记录sftp 日志

		err = lClient.SaveLog("dl-file", fullPath)
		if handleError(c, err) {
			return
		}
		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, f.Name()),
		}
		c.DataFromReader(http.StatusOK, fi.Size(), "application/octet-stream", f, extraHeaders)
		return
	} else if sourceType == "dir" {
		buf := new(bytes.Buffer)
		w := zip.NewWriter(buf)
		err := zipAddFiles(w, sftpClient, fullPath, "/")
		if handleError(c, err) {
			return
		}
		// Make sure to check the error on Close.
		err = w.Close()
		if handleError(c, err) {
			return
		}
		//记录sftp 日志
		err = lClient.SaveLog("dl-dir", fullPath)
		if handleError(c, err) {
			return
		}
		dName := time.Now().Format("2006_01_02T15_04_05Z07.zip")
		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, dName),
		}
		c.DataFromReader(http.StatusOK, int64(buf.Len()), "application/zip", buf, extraHeaders)
		return
	} else {
		jsonError(c, "query param type must be 'file' or 'dir'")
	}

}
func SftpCat(c *gin.Context) {
	fullPath := c.Query("path")
	lClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	defer lClient.Close()

	sftpClient := lClient.SftpClient
	fileInfo, err := sftpClient.Stat(fullPath)
	if handleError(c, err) {
		return
	}
	if fileInfo.IsDir() {
		jsonError(c, fullPath+" 是目录不能查查看文件内容")
		return
	}
	f, err := sftpClient.Open(fullPath)
	b, err := ioutil.ReadAll(f)
	if handleError(c, err) {
		return
	}
	//记录sftp 日志
	err = lClient.SaveLog("cate", fullPath)
	if handleError(c, err) {
		return
	}
	//c.String(200,"utf-8",file)
	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": string(b), "msg": fileInfo.Name()})
}
func zipAddFiles(w *zip.Writer, sftpC *sftp.Client, basePath, baseInZip string) error {
	// Open the Directory
	files, err := sftpC.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("sftp 读取目录 %s 失败:%s", basePath, err)
	}

	for _, file := range files {
		thisFilePath := basePath + "/" + file.Name()
		if file.IsDir() {

			err := zipAddFiles(w, sftpC, thisFilePath, baseInZip+file.Name()+"/")
			if err != nil {
				return fmt.Errorf("递归目录%s 失败:%s", thisFilePath, err)
			}
		} else {

			dat, err := sftpC.Open(thisFilePath)
			if err != nil {
				return fmt.Errorf("sftp 读取文件失败 %s:%s", thisFilePath, err)
			}
			// Add some files to the archive.
			zipElePath := baseInZip + file.Name()
			f, err := w.Create(zipElePath)
			if err != nil {
				return fmt.Errorf("写入zip writer header失败 %s:%s", zipElePath, err)
			}
			b, err := ioutil.ReadAll(dat)
			if err != nil {
				return fmt.Errorf("ioutil read all failed", err)
			}
			_, err = f.Write(b)
			if err != nil {
				return fmt.Errorf("写入zip writer 内容 bytes失败:%s", err)
			}
		}
	}
	return nil
}
