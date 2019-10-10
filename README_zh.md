# SSH-Fortress

## 1. What does it do?

1. Make your cluster servers be more safe by expose your SSH connection through SSH-Fortress server
2. Login your SSH server through the SSH-Fortress Web Interface and record all input and output history commands.
2. Manage your cluster server's SSH Account by SSH-Fortress with Web Account
3. Manage a server's files by SSH-Fortress's SFTP-web-interface
4. Easily login into your private Cluster by SSH Proxy provided by SSH-Fortress-Proxy


- 普通账号: 直接点击登陆(请输入验证码)


## 2.功能规划(一下功能都可以实现)

现阶段ssh/sfp都是使用的集群ssh账号,没有使用用户的ssh账号, 鉴权使用的是web账号和机器的关系.

后期要给每一个web账号在授权的机器上创建ssh账号,普通用户使用自己的ssh账号进行操作.

1. 完成: web账号授权给机器的时候同时,在目标机器创建web账号的ssh账号,以后ssh/sftp使用web用户独立账号登陆,
8. 完成: 获取管理机器的硬盘,内存,网络,CPU,操作,SSH账号...系统信息
4. 完成: 在外网/内网部署的情况下,支持使用github.com账号登陆20190805
3. todo: 用户登陆支持LDAP协议,
5. 完成: 在配置文件配置策略,增加定时任务,来清理日志日,和危险操作报警(邮件/sms/微信)
6. 完成: 配置文件增加策略,过滤一些 危险的ssh 命令,
7. todo: web界面列表支持字段搜索更具用户需求.
9. todo: 支持登陆到Kubernetes Pod 终端(terminal)
10. 自定义ssh协议, 可以在termninal中通过堡垒机命令行登陆目标机器.
fixed:修改不能修改/etc/sudoers 文件,permission 权限不够的问题

## 3.编译部署

特色:

1. 部署简单: go编译好的二进制文件(包含前端代码和后端代码,前端代码已经被打包), 一个`config.toml`配置文件
2. 支持多种SQL数据库,目前使用的是mysql5.7, 只需简单修就可以支持SQLite3/Postgres登陆SQL数据库

### Example 配置文件

`config.toml` app 会再 `/etc/sshfortress/config.toml` ,  `$HOME/.sshfortress/config.toml`,  `${app可执行文件目录}/config.toml`  三个路径中一个寻找配置文件

```toml
[app]
    name="frotress.mojotv.cn"
    addr=":3333"
    verbose= true
    jwt_expire=24# jwt 过期时间 单位小时
    secret="asdf4e8hcjvbkjclkjkklfgki843895iojfdnvufh98" #jwt
    ssh_log_days=7 #每台机器ssh终端日志保留天数
    sftp_log_days=7 #每台机器sftp日志保留天数
    signin_log_days=30 #登陆日志保留天数
    ssh_user_prefix="sshfortress_" #sshfortress 授权用户到机器的时候 会在目标机器上创建ssh账号, 前缀+ user.Name 就构成了 ssh账号用户名, 防止用户overwrite root这样的系统账号
[github] #github.com OAuth2 登陆地址
    client_id="d0b29360a033d0c4dc18"
    client_secret="89b272eeb22f373d8aa6c3986a8dbbc4edbfc64a"
    callback_url="http://sshfortress.mojotv.cn:3333/#/"
[db] #mysql数据地址5.7
    host = "sshfortress.mojotv.cn"
    user = "root"
    dbname = "sshfortress"
    password = "!Venom2018"
    port = 3306
```


### 编译
- 开发请使用go > 1.12版本 开启go mod特性
- [felixbin](felixbin) 是go转换的前端vuejs(编译之后)代码. 使用 `sshfortress ginbin -s ${前端编译之后的代码目录} -d ${输出的go代码包路径felixbin} ` 命令生成的代码不能手动修改
- 编译 `git clone https://git.corp.mojotv.cn/aio-cloud/sshfortress.git;cd frotress; go build; echo '也可以go install, 到 $GOBIN 目录'; `
- 编写配置文件到任意一个位置  `/etc/sshfortress/config.toml` ,`$HOME/.sshfortress`,`${app可执行文件目录}/config.toml`
- `cd #{编译好的目录}` 运行 `./sshfortress run`.  不需要关系前端代码. 前端编译之后代码已经包含再编译之后的二进制文件中.

### superviosr 进程守护
`/etc/supervisord/sshfortress.ini`

```ini
[program:sshfortress]
directory=/home/zhouqing1/sshfortress
command=/home/zhouqing1/sshfortress/sshfortress run
autostart=true
autorestart=true
startsecs=10
user=root
chmod=0777
numprocs=1
redirect_stderr=true
stdout_logfile=/home/zhouqing1/sshfortress/sp.log
```

## 怎么解决私有云集群访问的问题
私有云里面一台机器设置对sshfortress堡垒机可以访问的代理,这台机器代号SshProxyServer

```go
	ssh := &easyssh.MakeConfig{
		User:    "drone-scp",// 
		Server:  "target.private.cluster.server",//私有云代理服务器
		Port:    "22",
		KeyPath: "./tests/.ssh/id_rsa",
		Proxy: easyssh.DefaultConfig{
			User:    "drone-scp",//代理服务器ssh用户名
			Server:  "private.cluster.proxy.server",//私有云代理服务器 
			Port:    "22",//端口
			KeyPath: "./tests/.ssh/id_rsa",//代理云服务器密钥or 密码
		},
	}
```