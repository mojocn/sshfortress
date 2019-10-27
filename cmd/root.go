package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sshfortress",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if isShowVersion {
			fmt.Printf("Golang Env: %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
			fmt.Printf("UTC build time:%s\n", buildTime)
			fmt.Printf("Build from Github repo version:%s\n", gitHash)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(bTime, gHash string) {
	buildTime = bTime
	gitHash = gHash
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var buildTime, gitHash string
var verbose, isShowVersion bool

func init() {
	//cobra.OnInitialize(initFunc)
	rootCmd.Flags().BoolVarP(&isShowVersion, "version", "V", false, "打印编译信息")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
}

func initViperConfigFile(configFile string) error {
	viper.SetConfigName(configFile)           // name of config file (without extension)
	viper.AddConfigPath("/etc/sshfortress/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.sshfortress") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	return viper.ReadInConfig()               // Find and read the config file
}
