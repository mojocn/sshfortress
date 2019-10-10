FROM golang:1.13

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/" GO111MODULE=on


WORKDIR /go/src/sshfortress
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN sshfortress -V
# 需要使用 volume 映射config.toml配置文件到 WORKDIR/config.toml
# EXPOSE 更具配置文件暴露端口


#映射配置文件
VOLUME /go/src/sshfortress/config.toml

CMD ["sshfortress","run"]