package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sshfortress/model"
	"sshfortress/util"
	"strconv"
	"strings"
	"time"
)

func AnyWebTerminal(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if handleError(c, err) {
		return
	}
	defer wsConn.Close()

	cols, err := strconv.Atoi(c.DefaultQuery("cols", "120"))
	if wshandleErrorPro(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "32"))
	if wshandleErrorPro(wsConn, err) {
		return
	}

	sshAddr, exist := c.GetQuery("a")
	if !exist {
		wsConn.WriteControl(websocket.CloseMessage, []byte("ssh addr parameter is not exist"), time.Now())
		return
	}
	//save your server from brute force password
	if strings.Contains(sshAddr, "39.106.87.48") || strings.Contains(sshAddr, "127.0.0.1") || strings.Contains(sshAddr, "localhost") || strings.Contains(sshAddr, "::1") || strings.Contains(sshAddr, ".mojotv.cn") {
		logrus.WithField("ip", c.ClientIP()).Error("criminal criminal criminal, protect my own serve")
		return
	}

	sshPassword, exist := c.GetQuery("p")
	if !exist {
		wsConn.WriteControl(websocket.CloseMessage, []byte("ssh password parameter is not exist"), time.Now())
		return
	}

	sshUser, exist := c.GetQuery("u")
	if !exist {
		wsConn.WriteControl(websocket.CloseMessage, []byte("ssh user parameter is not exist"), time.Now())
		return
	}

	sshClient, err := util.CreateSimpleSshClient(sshUser, sshPassword, sshAddr)
	if wshandleErrorPro(wsConn, err) {
		return
	}
	defer sshClient.Close()

	sws, err := model.NewLogicSshWsSession(cols, rows, true, sshClient, wsConn, nil)
	if wshandleErrorPro(wsConn, err) {
		return
	}
	defer sws.Close()

	quitChan := make(chan bool, 3)
	sws.Start(quitChan)
	go sws.Wait(quitChan)

	<-quitChan

	//保存日志
	//logrus.Info("websocket finished")
}
