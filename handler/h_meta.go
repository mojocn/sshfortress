package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Meta(c *gin.Context) {
	data := gin.H{
		"github_client_id":    viper.GetString("github.client_id"),
		"github_callback_url": viper.GetString("github.callback_url"),
		"grafana_base_url":    viper.GetString("grafana.base_url"),
	}
	jsonData(c, data)
}
