package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"sshfortress/cronjob"
	"sshfortress/felixbin"
	"sshfortress/handler"
	"sshfortress/model"
	"time"
)

// appCmd represents the sshw command
var appCmd = &cobra.Command{
	Use:   "run",
	Short: "运行sshfortress 项目",
	Long: `the demo website is http://sshfortress.mojotv.cn:3333
如果没有配置文件就创建SQLite数据库
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initViperConfigFile(configFile); err != nil {
			logrus.WithError(err).Fatal("加载配置文件失败")
		}
		model.AppSecret = viper.GetString("app.secret")
		model.ExpireTime = time.Hour * time.Duration(viper.GetInt("app.jwt_expire"))
		model.AppIss = viper.GetString("app.name")
		err := model.CreateMysqlDb(verbose)
		if err != nil {
			logrus.WithError(err).Fatal("数据库加载失败")
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
var configFile, addr string
var isMigrate bool

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().StringVarP(&configFile, "config", "c", "config", "指定加载的配置文件")
	appCmd.Flags().StringVarP(&addr, "listen", "l", ":8360", "app监听的端口")
	appCmd.Flags().BoolVarP(&isMigrate, "migrate", "m", true, "应用启动是否执行数据库迁移默认TRUE")
}

func runApiAndH5() error {
	//config jwt variables
	if !verbose {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	//开发接口
	{
		//chromedp run screen shot
		r.GET("open/chromedp/shot", handler.ChromedpShot)
	}

	r.GET("metrics", gin.WrapH(promhttp.Handler()))
	r.MaxMultipartMemory = 1024 * 1024 * 50 //100M 上传文件路径 todo::配置文件中进行配置
	//sever static file in http's root path
	binStaticMiddleware, err := felixbin.NewGinStaticBinMiddleware("/")
	if err != nil {
		return err
	}
	r.Use(binStaticMiddleware)
	api := r.Group("api").Use(handler.MwPrometheusHttp)
	//监控检查给docker health check 使用
	api.GET("health-check", func(context *gin.Context) {
		context.AbortWithStatusJSON(200, gin.H{"ok": true})
	})
	//不需要token
	{
		r.POST("api/login", handler.MwCaptchaCheck(), handler.Login)
		r.GET("api/login-github", handler.LoginGithub)
		r.POST("api/register", handler.UserCreate)
		r.GET("api/captcha", handler.GetCaptcha)
		r.GET("api/meta", handler.Meta)
	}
	api.GET("ws-any-term", handler.AnyWebTerminal)

	mwAdmin := handler.MwUserRole(2, "需要管理员权限")
	authG := api.Use(handler.JwtMiddleware)
	{
		//websocket ssh
		authG.GET("ws/:id", handler.MachineWsSshTerm)
		{
			authG.GET("cluster-ssh", mwAdmin, handler.ClusterSshAll)
			authG.POST("cluster-ssh", mwAdmin, handler.ClusterSshCreate)
			authG.PATCH("cluster-ssh", mwAdmin, handler.ClusterSshUpdate)
			authG.DELETE("cluster-ssh/:id", mwAdmin, handler.ClusterSshDelete)
			authG.GET("cluster-ssh/:id", handler.ClusterSshOne)
			authG.POST("cluster-ssh-bind", mwAdmin, handler.ClusterSshBindMachines)
		}

		{
			// config table in SQL database
			authG.GET("config", mwAdmin, handler.ConfigAll)
			authG.PATCH("config", mwAdmin, handler.ConfigUpdate)
		}

		{
			authG.GET("cluster-jumper", mwAdmin, handler.ClusterJumperAll)
			authG.POST("cluster-jumper", mwAdmin, handler.ClusterJumperCreate)
			authG.PATCH("cluster-jumper", mwAdmin, handler.ClusterJumperUpdate)
			authG.DELETE("cluster-jumper/:id", mwAdmin, handler.ClusterJumperDelete)
			authG.GET("cluster-jumper/:id", handler.ClusterJumperOne)
			authG.POST("cluster-jumper-bind", mwAdmin, handler.ClusterJumperBindMachines)
		}
		//用户组权限控制
		{
			authG.GET("filter-group", mwAdmin, handler.SshFilterGroupAll)
			authG.POST("filter-group", mwAdmin, handler.SshFilterGroupCreate)
			authG.PATCH("filter-group", mwAdmin, handler.SshFilterGroupUpdate)
			authG.DELETE("filter-group/:id", mwAdmin, handler.SshFilterGroupDelete)
			authG.GET("filter-group/:id", handler.SshFilterGroupOne)
		}
		{
			authG.GET("machine", handler.MachineAll)
			authG.POST("machine", handler.MachineCreate)
			authG.PATCH("machine", handler.MachineUpdate)
			authG.DELETE("machine/:id", handler.MachineDelete)
			authG.GET("machine/:id", handler.MachineOne)
			authG.GET("machine/:id/hardware", handler.MachineHardware)
		}
		{
			authG.GET("ssh-log", handler.SshLogAll)
			authG.PATCH("ssh-log", handler.SshLogUpdate)
			authG.POST("ssh-log-rm", handler.SshLogDelete)
		}
		{

			authG.GET("sftp-log", handler.SftpLogAll)
			authG.PATCH("sftp-log", handler.SftpLogUpdate)
			authG.POST("sftp-log-rm", handler.SftpLogDelete)
		}

		{
			authG.GET("feedback", mwAdmin, handler.FeedbackAll)
			authG.POST("feedback", handler.FeedbackCreate)
			authG.PATCH("feedback", mwAdmin, handler.FeedbackUpdate)
			authG.DELETE("feedback-rm", mwAdmin, handler.FeedbackDelete)
		}

		{
			authG.GET("user", mwAdmin, handler.UserAll)
			authG.POST("user", mwAdmin, handler.UserCreate)
			authG.DELETE("user/:id", mwAdmin, handler.UserDelete)
			authG.GET("user/:id", mwAdmin, handler.UserOne)
			authG.PATCH("user", handler.UserUpdate)
		}
		{
			authG.GET("sftp/:id", handler.SftpLs)
			authG.GET("sftp/:id/dl", handler.SftpDl)
			authG.GET("sftp/:id/cat", handler.SftpCat)
			authG.GET("sftp/:id/rm", handler.SftpRm)
			authG.GET("sftp/:id/rename", handler.SftpRename)
			authG.GET("sftp/:id/mkdir", handler.SftpMkdir)
			authG.POST("sftp/:id/up", handler.SftpUp)
		}
		{
			authG.GET("machine-user/user-ids", mwAdmin, handler.MachineUserUserIds)
			authG.GET("machine-user/machine-ids", mwAdmin, handler.MachineUserMachineIds)
			authG.POST("machine-user/bind-machine", mwAdmin, handler.MachineUserBindMachines)
			authG.POST("machine-user/bind-user", mwAdmin, handler.MachineUserBindUsers)
		}
		{
			authG.GET("signin-log", mwAdmin, handler.SigninLogAll)
			authG.POST("signin-log-rm", mwAdmin, handler.SigninLogDelete)
		}
	}
	if err := r.Run(addr); err != nil {
		return err
	}
	return nil
}
