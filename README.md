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
DOMAINS = {"domain.cn":"http://127.0.0.1:8080","api.domain.cn":"http://192.168.1.100"}
```

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
DOMAINS = {"domain.cn":"http://127.0.0.1:8080","api.domain.cn":"http://192.168.1.100"}
```

### 压力测试
与frp(https://github.com/fatedier/frp )进行压力测试对比，分别部署内网穿透服务端及客服端后，对内网服务器的某个http站点进行压测。

压测工具ab，压测页面是一个大小为3.2kb的html，这里只贴上Requests per second这项数据。

并发100，请示数1000；

```
frp：
Requests per second:    114.48 [#/sec] (mean)
Requests per second:    131.48 [#/sec] (mean)
Requests per second:    106.38 [#/sec] (mean)

proxy(v0.2.0):
Requests per second:    135.92 [#/sec] (mean)
Requests per second:    122.72 [#/sec] (mean)
Requests per second:    117.42 [#/sec] (mean)

proxy(v0.1.0)
Requests per second:    26.24 [#/sec] (mean)
Requests per second:    28.53 [#/sec] (mean)
Requests per second:    21.98 [#/sec] (mean)
```

并发100，请示数1000；
```
frp:
Requests per second:    131.19 [#/sec] (mean)
Requests per second:    128.39 [#/sec] (mean)
Requests per second:    111.74 [#/sec] (mean)

proxy:
Requests per second:    95.16 [#/sec] (mean)
Requests per second:    89.68 [#/sec] (mean)
Requests per second:    121.06 [#/sec] (mean)
```
与v0.1.0相比，使用gRPC后性能明显提升，测试结果甚至与frp相当，不过当我逐渐提升并发量进行测试时,程序就不稳定甚至宕机，frp则依然稳定执行。
