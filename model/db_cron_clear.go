package model

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

const timef = "2006-01-02 15:04:05"

func ClearLogTable() (err error) {
	if db == nil {
		err = errors.New("请先出初始化db")
		return
	}
	signinLogM := SigninLog{}
	if db.HasTable(&signinLogM) {
		d := viper.GetInt("app.signin_log_days")
		if d < 1 {
			d = 7
		}
		t := time.Now().Add(-time.Hour * time.Duration(d*24))
		err = db.Unscoped().Where(`created_at < ?`, t.Format(timef)).Delete(signinLogM).Error
		if err != nil {
			return
		}
	}

	sshLogM := SshLog{}
	if db.HasTable(&sshLogM) {
		d := viper.GetInt("app.ssh_log_days")
		if d < 1 {
			d = 7
		}
		t := time.Now().Add(-time.Hour * time.Duration(d*24))
		err = db.Unscoped().Where(`created_at < ?`, t.Format(timef)).Delete(sshLogM).Error
		if err != nil {
			return
		}
	}

	sftpLogM := SftpLog{}
	if db.HasTable(&sftpLogM) {
		d := viper.GetInt("app.sftp_log_days")
		if d < 1 {
			d = 7
		}
		t := time.Now().Add(-time.Hour * time.Duration(d*24))
		err = db.Unscoped().Where(`created_at < ?`, t.Format(timef)).Delete(sftpLogM).Error
		if err != nil {
			return
		}
	}
	return nil
}
