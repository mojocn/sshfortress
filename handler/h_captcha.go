package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/util"
)

func GetCaptcha(c *gin.Context) {
	id, imageString := util.GetCaptchaImage()
	jsonData(c, gin.H{"id": id, "image": imageString})
}
