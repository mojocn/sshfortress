package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sshfortress/model"
	"strings"
)

const jwtCtxUidKey = "authedUserId"
const bearerLength = len("Bearer ")

func JwtMiddleware(c *gin.Context) {
	token, ok := c.GetQuery("_t")
	if !ok {
		hToken := c.GetHeader("Authorization")
		if len(hToken) < bearerLength {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
			return
		}
		token = strings.TrimSpace(hToken[bearerLength:])
	}
	userId, err := model.JwtParseUser(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": err.Error()})
		return
	}
	//store the user Model in the context
	c.Set(jwtCtxUidKey, userId)
	c.Next()
	// after request
}
