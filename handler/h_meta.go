package handler

import (
	"github.com/gin-gonic/gin"
	"sshfortress/model"
)

func Meta(c *gin.Context) {
	data := gin.H{
		"github_client_id":    model.GithubClientId,
		"github_callback_url": model.GithubClientCallbackUrl,
	}
	jsonData(c, data)
}
