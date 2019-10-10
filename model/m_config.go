package model

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sshfortress/util"
	"time"
)

var AppSecret = "J%df4e8hcjvbkjclkjkklfgki843895iojfdnvufh98"
var AppIss = "sshfortress"
var ExpireTime = time.Hour * 24 * 30

var GithubClientId = ""
var GithubClientSecret = ""
var GithubClientCallbackUrl = ""

type ConfigQ struct {
	PaginationQ
	Config
}

//保存一些配置的值
type Config struct {
	BaseModel
	Key    string `gorm:"index" json:"key" form:"key"`
	Value  string `json:"value" form:"value"`
	Remark string `json:"remark" form:"remark"`
}

//All
func (m Config) All(q *PaginationQ) (list *[]Config, total uint, err error) {
	list = &[]Config{}
	tx := db.Model(m) //.Where("ancestor_path like ?", m.qAncetorPath())
	if m.Remark != "" {
		tx = tx.Where("remark like ?", "%"+m.Remark+"%")
	}
	if m.Key != "" {
		tx = tx.Where("`key` like ?", "%"+m.Key+"%")
	}

	if m.Value != "" {
		tx = tx.Where("value like ?", "%"+m.Value+"%")
	}
	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *Config) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	migrateOrLoadConfig()
	return db.Model(m).Update(m).Error
}

func insertOrCreateConfigItem(key, def, remark string) string {
	cfg := Config{Key: key}
	if db.Where("`key` = ?", key).First(&cfg).RecordNotFound() {
		cfg.Value = def
		cfg.Remark = remark
		if err := db.Create(&cfg).Error; err != nil {
			logrus.WithError(err).WithField("key", key).WithField("v", def).Error("写入config初始化值失败")
		}
	}
	return cfg.Value
}

func migrateOrLoadConfig() {
	//初始化 配置值 for jwt
	AppSecret = insertOrCreateConfigItem("app_secret", util.RandomDigitAndLetters(24), "JWT secret at least 24 byte")
	AppIss = insertOrCreateConfigItem("app_name", "ssh_fortress", "JWT Issuer Name")
	GithubClientId = insertOrCreateConfigItem("github.client_id", "d0b29360a033d0c4dc18", "github OAuth2 Client Id")
	GithubClientSecret = insertOrCreateConfigItem("github.client_secret", "89b272eeb22f373d8aa6c3986a8dbbc4edbfc64a", "github OAuth2 Client Secret")
	GithubClientCallbackUrl = insertOrCreateConfigItem("github.callback_url", "https://sshfortress.mojotv.cn/#/", "github OAuth2 Callback URL")
}
