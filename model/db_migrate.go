package model

import (
	"github.com/sirupsen/logrus"
	"time"
)

func RunMigrate() error {
	ms := []interface{}{User{}, ClusterSsh{}, Machine{}, MachineUser{}, Config{}, SshLog{}, SigninLog{}, SftpLog{}, Feedback{}, ClusterJumper{}, SshFilterGroup{}}
	for idx, v := range ms {
		if err := db.AutoMigrate(v).Error; err != nil {
			logrus.WithError(err).Error("迁移模型失败:", idx)
		}
	}

	god := User{
		RealName:         "SuperAdmin",
		Email:            "admin@sshfortress.cn",
		InputPassword:    "admin",
		Mobile:           "13312345678",
		Role:             2,
		InputSshPassword: "admin",
		Name:             "hydra"}
	god.CreatedAt = time.Now()
	god.UpdatedAt = time.Now()
	ex := time.Now().Add(time.Hour * 24 * 365 * 10)
	god.ExpiredAt = &ex
	god.makePassword()
	err := db.FirstOrCreate(&god, User{Email: god.Email}).Error
	if err != nil {
		logrus.WithError(err).Error("创建初始化用户失败")
	}

	migrateOrLoadConfig()
	return nil
}
