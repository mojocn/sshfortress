package cmd

import (
	"github.com/spf13/cobra"
	"sshfortress/fssh"
)

var openApiCmd = &cobra.Command{
	Use:   "oi",
	Short: "run open api http-server",
	Long:  `a set of open api server`,
	Run: func(cmd *cobra.Command, args []string) {
		fssh.Run()
	},
}

func init() {
	rootCmd.AddCommand(openApiCmd)
}
