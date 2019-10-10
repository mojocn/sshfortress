package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"sshfortress/cronjob"
	"sshfortress/model"
)

var sqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "run sshfortress with sqlite",
	Long:  `run sshfortress is much more simple without any config file and mysql-connection, it will generate a sqlite database file atomically`,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.CreateSqliteDb(verbose)
		if err != nil {
			logrus.WithError(err).Fatal("create/open SQLite3 file failed")
		}
		defer model.Close()
		go cronjob.RunsshfortressCron()
		//执行数据库迁移
		if isMigrate {
			err = model.RunMigrate()
			if err != nil {
				log.Fatal(err)
			}
		}
		//加载 ssh 命令配置策略到内存中从数据库
		err = runApiAndH5()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sqliteCmd)
	sqliteCmd.Flags().StringVarP(&addr, "listen", "l", ":8360", "app监听的端口")
}
