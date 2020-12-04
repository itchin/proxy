通过公网服务器代理转发请求，使内网http服务能被外网访问到，适合于接口开发测试使用。


### 二进制下载

https://github.com/itchin/proxy/releases


### 部署说明

需先安装golang开发环境，并设置服务端及客户端配置

```
git clone https://github.com/itchin/proxy

#下载所需组件
go mod tidy

#运行服务端
go run runServer.go

#运行客户端
go run runClient.go
```

### 服务端配置

复制 server.exam.ini 为 server.ini

```
CONSOLE_LOG = false
GRPC_HOST = :9090
HTTP_HOST = 0.0.0.0:9097
```

### 客户端配置

复制 client.exam.ini 为 client.ini

```
CONSOLE_LOG = false
GRPC_HOST = domain.cn:9090
GZIP_COMPRESSION = 5
DOMAINS = {"www.domain.cn":"http://127.0.0.1:8080","api.domain.cn":"http://192.168.1.100"}
```
