package fssh

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sshfortress/util"
)

func Run() {
	hostKeySigner, err := createOrLoadKeySigner()
	if err != nil {
		log.Fatal(err)
	}
	s := &ssh.Server{
		Addr:    ":88",
		Handler: homeHandler, //
		//PasswordHandler: passwordHandler,   不需要密码验证
	}
	s.AddHostKey(hostKeySigner)
	log.Fatal(s.ListenAndServe())
}

func passwordHandler(ctx ssh.Context, password string) bool {
	//check password and username
	//user := ctx.User()
	// 可以结合DB数据库定义用户验证用户登陆
	return true
	//return model.FsshUserAuth(user, password)
}

func homeHandler(s ssh.Session) {
	//tty 控制码打印彩色文字
	//mojotv.cn/tutorial/golang-term-tty-pty-vt100
	io.WriteString(s, fmt.Sprintf("\x1b[31;47mmojotv.cn sshfortress 堡垒机 自定义SSH, 当前登陆用户名: %s\x1b[0m\n", s.User()))

	ptyReq, winCh, isPty := s.Pty()
	if !isPty {
		io.WriteString(s, "不是PTY请求.\n")
		s.Exit(1)
		return
	}
	sshConf, err := util.NewSshClientConfig("test007", "test007", "password", "", "")
	if err != nil {
		io.WriteString(s, err.Error())
		s.Exit(1)
		return
	}
	//连接远程服务器SSH
	conn, err := gossh.Dial("tcp", "home.mojotv.cn:22", sshConf)
	if err != nil {
		io.WriteString(s, "unable to connect: "+err.Error())
		s.Exit(1)
		return
	}
	defer conn.Close()
	// 创建远程ssh session
	fss, err := conn.NewSession()
	if err != nil {
		io.WriteString(s, "unable to create fss: "+err.Error())
		s.Exit(1)
		return
	}
	defer fss.Close()

	// 配置terminal
	modes := gossh.TerminalModes{
		gossh.ECHO:          1,     // disable echoing
		gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// 请求为终端
	if err := fss.RequestPty(ptyReq.Term, ptyReq.Window.Height, ptyReq.Window.Width, modes); err != nil {
		io.WriteString(s, "request for pseudo terminal failed: "+err.Error())
		s.Exit(1)
		return
	}
	//监听终端size window 变化
	go func() {
		for win := range winCh {
			err := fss.WindowChange(win.Height, win.Width)
			if err != nil {
				io.WriteString(s, "windows size changed: "+err.Error())
				s.Exit(1)
				return
			}
		}
	}()

	//linux 一切接文件 io, 连接stdin stdout stderr
	//连接为终端到server
	fss.Stderr = s
	fss.Stdin = s
	fss.Stdout = s
	if err := fss.Shell(); err != nil {
		io.WriteString(s, "failed to start shell: "+err.Error())
		s.Exit(1)
		return
	}
	fss.Wait()
}

//创建key 来验证 host public
func createOrLoadKeySigner() (gossh.Signer, error) {
	keyPath := filepath.Join(os.TempDir(), "fssh.rsa")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
		stderr, err := exec.Command("ssh-keygen", "-f", keyPath, "-t", "rsa", "-N", "").CombinedOutput()
		output := string(stderr)
		if err != nil {
			return nil, fmt.Errorf("Fail to generate private key: %v - %s", err, output)
		}
	}
	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return gossh.ParsePrivateKey(privateBytes)
}
