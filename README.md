通过公网服务器代理转发请求，使内网http服务能被外网访问到，适合于接口开发测试使用。


### 二进制下载

https://github.com/itchin/proxy/releases


### 部署说明

需先安装golang开发环境，并设置服务端及客户端配置

```
git clone https://github.com/itchin/proxy

#下载所需组件
go mod vendor

#运行服务端
go run server_ctl.go

#运行客户端
go run client_ctl.go
```

### 服务端配置

复制 server.exam.ini 为 server.ini

```
# HTTP服务器最大连接数
MAX_CONN = 1024
# 最大活跃协程数量
MAX_ACTIVE = 20
# http请求过期时间(秒)
HTTP_TIMEOUT = 30
CONSOLE_LOG = false
HTTP_HOST = 0.0.0.0:80
GRPC_HOST = :9090
```

### 客户端配置

复制 client.exam.ini 为 client.ini

```
WORKERS = 5
GRPC_HOST = domain.cn:9090
HTTP_TIMEOUT = 30
CONSOLE_LOG = false
# GZIP压缩，0~9，0为关闭
GZIP_COMPRESSION = 0
DOMAINS = {"domain.cn":"http://127.0.0.1:8080","api.domain.cn":"http://192.168.1.100"}
HEARTBEAT = 0
```
