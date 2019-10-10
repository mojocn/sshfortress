package handler

import (
	"github.com/gin-gonic/gin"
)

func MwUserRole(role uint, msg string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := mwJwtUser(c)
		if err != nil {
			jsonError(c, err)
			return
		}
		if u.Role != role {
			jsonError(c, msg)
			return
		}
		c.Next()
	}
}
