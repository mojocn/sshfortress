# SSH-Fortress

## 1. What does it do?

1. Make your cluster servers be more safe by expose your SSH connection through SSH-Fortress server
2. Login your SSH server through the SSH-Fortress Web Interface and record all input and output history commands.
2. Manage your cluster server's SSH Account by SSH-Fortress with Web Account
3. Manage a server's files by SSH-Fortress's SFTP-web-interface
4. Easily login into your private Cluster by SSH Proxy provided by SSH-Fortress-Proxy


## 2. build and run
```bash
git clone https://github.com/mojocn/sshfortress.git && cd sshfortress;
go build
echo "run the app with SQLite database"
./sshfortress sqlite -v --listen=':3333'
echo "run the app with Mysql database, you need a config.toml file in your sshfortress binary folder"
./sshfortress run -v --listen=':3333'

```
Docker pull `docker pull mojotvcn/sshfortress`

### 2.1 config.toml
The config.toml file should in sshfortress binary folder.  config.toml works with command `sshfortress run`. Command `sshfortress sqlite` can run with the config file.

```toml
[app]
    name="frotress.mojotv.cn"
    addr=":8360"
    verbose= true
    jwt_expire=240 #hour
    secret="asdf4e8hcjvbkjclkjkklfgki843895iojfdnvufh98" #jwt secret
[db]
    # mysql database connection
    host = "127.0.0.1"
    user = "root"
    dbname = "sshfortress"
    password = "your_mysql_password"
    port = 3306

[github] #github.com OAuth2
    client_id="d0b29360a088d0c4dc18"
    client_secret="89b272eeb22f373d8aa688986a8dbbc4edbfc64a"
    callback_url="http://sshfortress.mojotv.cn/#/"
```
## 3. Online demo

[https://sshfortress.mojotv.cn/#/login](https://sshfortress.mojotv.cn/#/login)

just click the login button, the default password has input for you, user `admin@sshfortress.cn` password: `admin`,

### 3.1 Universal Web SST Terminal

- `URL` : `https://sshfortress.mojotv.cn/#/any-term`  eg: `https://sshfortress.mojotv.cn/#/any-term?a=home.mojotv.cn&p=test007&u=test007&z=1`
- URL-ARG  `a` : SSH Address with Port eg: `home.mojotv.cn` `home.mojotv.cn:22`
- URL-ARG  `u` : SSH Username eg: `test007`
- URL-ARG  `p` : SSH Password eg: `test007`
- URL-ARG  `z` : Not Use Zend Mode eg: `1`


## 4. Run With supervisor & nginx

sshfortress.mojotv.cn.conf
```bash
server {
        server_name sshfortress.mojotv.cn;
        charset utf-8;
        location /api/ws-any-term
        {
                proxy_pass http://127.0.0.1:8360;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "Upgrade";
                proxy_set_header X-Real-IP $remote_addr;
         }

        location /api/ws/
        {
                proxy_pass http://127.0.0.1:8360;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "Upgrade";
                proxy_set_header X-Real-IP $remote_addr;
         }
        location / {
           proxy_set_header X-Forwarded-For $remote_addr;
           proxy_set_header Host $http_host;
           proxy_pass http://127.0.0.1:8360;
        }
        access_log  /data/wwwlogs/sshfortress.mojotv.cn.log;


    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/sshfortress.mojotv.cn/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/sshfortress.mojotv.cn/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot
}
```

Supervisor config file: `sshfortress.ini`
```ini
[program:sshfortress.mojotv.cn]
command=/data/sshfortress/bin/sshfortress sqlite
autostart=true
autorestart=true
startsecs=10
user=root
chmod=0777
numprocs=1
redirect_stderr=true
stdout_logfile=/data/sshfortress/supervisor.log
```

## 5. Reference

1. [idea from my another repo: libragen/felix](https://github.com/libragen/felix)
2. [How to run SSH-Terminal in browser](https://mojotv.cn/2019/05/27/xtermjs-go)
3. [Dockerhub image](https://hub.docker.com/r/mojotvcn/sshfortress/dockerfile)