package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sshfortress/model"
	"sshfortress/stat"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func MachineWsSshTerm(c *gin.Context) {
	stat.GaugeVecApiMethod.WithLabelValues("WS").Inc()
	start := time.Now()
	defer func() {
		end := time.Now()
		d := end.Sub(start) / time.Millisecond
		stat.GaugeVecApiDuration.WithLabelValues("WS").Set(float64(d))
	}()
	user, err := mwJwtUser(c)
	if wshandleError(err) {
		return
	}
	cols, err := strconv.Atoi(c.DefaultQuery("cols", "120"))
	if wshandleError(err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "32"))
	if wshandleError(err) {
		return
	}
	idx, err := parseParamID(c)
	if wshandleError(err) {
		return
	}

	logicClient, err := model.CreateLogicSshClient(user, idx, c.ClientIP())
	if wshandleError(err) {
		return
	}
	defer logicClient.Close()

	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if handleError(c, err) {
		return
	}
	defer wsConn.Close()

	sfg := user.MustSshFilterGroup()
	sws, err := model.NewLogicSshWsSession(cols, rows, user.IsAdmin(), logicClient.SshClient, wsConn, sfg)
	if wshandleError(err) {
		return
	}
	defer sws.Close()

	quitChan := make(chan bool, 3)
	sws.Start(quitChan)
	go sws.Wait(quitChan)

	<-quitChan

	//保存日志
	err = logicClient.SaveLog(sws.IsFlagged, sws.HasEditor, sws.Tlog)
	if wshandleError(err) {
		return
	}
	//logrus.Info("websocket finished")
}
