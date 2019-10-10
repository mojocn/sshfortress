package cmd

import (
	"github.com/spf13/cobra"
	"sshfortress/fssh"
)

var fsshCmd = &cobra.Command{
	Use:   "fssh",
	Short: "运行sshfortress自定义的ssh:fssh服务",
	Long:  `默认端口:7777`,
	Run: func(cmd *cobra.Command, args []string) {
		fssh.Run()
	},
}

func init() {
	rootCmd.AddCommand(fsshCmd)
}
