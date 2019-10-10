package cronjob

import (
	"github.com/sirupsen/logrus"
	"sshfortress/model"
)

func RunsshfortressCron() {
	s := NewScheduler()
	s.Every(24).Hours().Do(doClearLogJob)
	<-s.Start()
}

func doClearLogJob() {
	if err := model.ClearLogTable(); err != nil {
		logrus.WithError(err).Error("cron job run error")
	}
}
