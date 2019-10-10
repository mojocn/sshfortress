package model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

type safeBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *safeBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}
func (w *safeBuffer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Bytes()
}
func (w *safeBuffer) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Reset()
}

const (
	wsMsgCmd    = "cmd"
	wsMsgResize = "resize"
)

type wsMsg struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

type LogicSshWsSession struct {
	stdinPipe       io.WriteCloser
	comboOutput     *safeBuffer //ssh 终端混合输出
	inputFilterBuff *safeBuffer //用来过滤输入的命令和ssh_filter配置对比的
	session         *ssh.Session
	wsConn          *websocket.Conn
	isAdmin         bool
	IsFlagged       bool `comment:"当前session是否包含禁止命令"`

	sshFilters JsonArraySshFilter

	Tlog      JsonArrayString
	HasEditor bool
}

func NewLogicSshWsSession(cols, rows int, isAdmin bool, sshClient *ssh.Client, wsConn *websocket.Conn, sfg *SshFilterGroup) (*LogicSshWsSession, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(safeBuffer)
	inputBuf := new(safeBuffer)
	//ssh.stdout and stderr will write output into comboWriter
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	// Start remote shell
	if err := sshSession.Shell(); err != nil {
		return nil, err
	}
	sfs := JsonArraySshFilter{}
	if sfg != nil {
		sfs = sfg.Filters
	}
	return &LogicSshWsSession{
		stdinPipe:       stdinP,
		comboOutput:     comboWriter,
		inputFilterBuff: inputBuf,
		session:         sshSession,
		wsConn:          wsConn,
		isAdmin:         isAdmin,
		IsFlagged:       false,
		sshFilters:      sfs,
		Tlog:            JsonArrayString{},
	}, nil
}

//Close 关闭
func (sws *LogicSshWsSession) Close() {
	if sws.session != nil {
		sws.session.Close()
	}

	if sws.comboOutput != nil {
		sws.comboOutput = nil
	}
}
func (sws *LogicSshWsSession) Start(quitChan chan bool) {
	go sws.receiveWsMsg(quitChan)
	go sws.sendComboOutput(quitChan)
}

//receiveWsMsg  receive websocket msg do some handling then write into ssh.session.stdin
func (sws *LogicSshWsSession) receiveWsMsg(exitCh chan bool) {
	wsConn := sws.wsConn
	//tells other go routine quit
	defer setQuit(exitCh)
	for {
		select {
		case <-exitCh:
			return
		default:
			//read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				logrus.WithError(err).Error("reading webSocket message failed")
				return
			}
			//unmashal bytes into struct
			msgObj := wsMsg{}
			if err := json.Unmarshal(wsData, &msgObj); err != nil {
				logrus.WithError(err).WithField("wsData", string(wsData)).Error("unmarshal websocket message failed")
			}
			switch msgObj.Type {
			case wsMsgResize:
				//handle xterm.js size change
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					if err := sws.session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
						logrus.WithError(err).Error("ssh pty change windows size failed")
					}
				}
			case wsMsgCmd:
				//handle xterm.js stdin
				decodeBytes, err := base64.StdEncoding.DecodeString(msgObj.Cmd)
				if err != nil {
					logrus.WithError(err).Error("websock cmd string base64 decoding failed")
				}
				sws.sendWebsocketInputCommandToSshSessionStdinPipe(decodeBytes)
			}
		}
	}
}

//sendWebsocketInputCommandToSshSessionStdinPipe
func (sws *LogicSshWsSession) sendWebsocketInputCommandToSshSessionStdinPipe(cmdBytes []byte) {
	//保存整行input
	var lineCommand []byte
	for _, bb := range cmdBytes {
		//判断命令是否开始换行或者;
		if bb == '\r' || bb == ';' || bb == '\n' {
			lineCommand = sws.inputFilterBuff.Bytes()
			sws.inputFilterBuff.Reset()
			//匹配配置的命令策略
		} else {
			_, err := sws.inputFilterBuff.Write([]byte{bb})
			if err != nil {
				logrus.WithError(err).Error("sws.inputFilterBuff.Write")
			}
		}
	}
	//匹配文本编辑器
	if len(lineCommand) > 0 {
		isEditor, err := regexp.Match(`\b(vim|vi|nano|emacs|gedit|kate|kedit)\b`, lineCommand)
		if err != nil {
			logrus.WithError(err).Error("检测文本编辑器失败")
		}
		if isEditor {
			sws.HasEditor = true
		}
	}

	//普通用户正则检测命令
	if !sws.isAdmin && len(lineCommand) > 0 {
		//处理命令过滤的问题;
		//匹配配置的命令策略
		for _, rule := range sws.sshFilters {
			patern := fmt.Sprintf(`\b%s\b`, strings.Trim(rule.Command, `\b`))
			isMatch, err := regexp.Match(patern, lineCommand)
			if err != nil {
				logrus.WithError(err).Error("regexp.Match(patern,rawCmdB)")
			}
			if isMatch {
				sws.IsFlagged = true
				//write warning msg into websocket terminal
				warning := fmt.Sprintf("\n\r \033[0;31m%s\033[0m\r\n", rule.Msg)
				sws.wsConn.WriteMessage(websocket.TextMessage, []byte(warning))
				//https://unicodelookup.com/#ctrl
				//制造clear命令
				cmdBytes = []byte{byte(025)}
				//return
			}
		}
	}
	if _, err := sws.stdinPipe.Write(cmdBytes); err != nil {
		logrus.WithError(err).Error("ws cmd bytes write to ssh.stdin pipe failed")
	}
}

func (sws *LogicSshWsSession) sendComboOutput(exitCh chan bool) {
	wsConn := sws.wsConn
	//todo 优化成一个方法
	//tells other go routine quit
	defer setQuit(exitCh)

	//every 120ms write combine output bytes into websocket response
	tick := time.NewTicker(time.Millisecond * time.Duration(10))
	//for range time.Tick(120 * time.Millisecond){}
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if sws.comboOutput == nil {
				return
			}
			bs := sws.comboOutput.Bytes()
			if len(bs) > 0 {
				err := wsConn.WriteMessage(websocket.TextMessage, bs)
				if err != nil {
					logrus.WithError(err).Error("ssh sending combo output to webSocket failed")
				}
				sws.comboOutput.buffer.Reset()
			}
			sws.writeLog(bs)

		case <-exitCh:
			return
		}
	}
}

func (sws *LogicSshWsSession) Wait(quitChan chan bool) {
	if err := sws.session.Wait(); err != nil {
		logrus.WithError(err).Error("ssh session wait failed")
		setQuit(quitChan)
	}
}

func (sws *LogicSshWsSession) writeLog(bs []byte) {
	if len(bs) > 0 {
		sws.Tlog = append(sws.Tlog, string(bs))
	}
}

func setQuit(ch chan bool) {
	ch <- true
}
