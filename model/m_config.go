package model

import (
	"github.com/sirupsen/logrus"
)

//保存一些配置的值
type Config struct {
	BaseModel
	Key    string `gorm:"index" json:"key" form:"key"`
	Value  string `json:"value" form:"value"`
	Remark string `json:"remark" form:"remark"`
}

func configKeyDefault(key, def string) string {
	cfg := Config{Key: key}
	if db.Where("key = ?", key).First(&cfg).RowsAffected != 1 {
		cfg.Value = def
		if err := db.Create(&cfg).Error; err != nil {
			logrus.WithError(err).WithField("key", key).WithField("v", def).Error("写入config初始化值失败")
		}
	}
	return cfg.Value
}
