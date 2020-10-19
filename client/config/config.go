package config

import (
    "fmt"

    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
    jsoniter "github.com/json-iterator/go"
    "gopkg.in/ini.v1"
)

var (
    // 是否在控制台中打印日志
    CONSOLE_LOG = true

    GRPC_HOST string

    // "远程域名": "本地协议://本地地址:端口"
    // 如：map[string]string{"api.me": "http://127.0.0.1:9000"}
    DOMAINS map[string]string

    // gzip压缩级别
    GZIP_COMPRESSION = 5

    // 心跳包间隔时间(秒)
    HEARTBEAT = 0

    // HTTP超时(秒)
    HTTP_TIMEOUT = 30
)

func init() {
    cfg, err := ini.Load("client.ini")
    if err != nil {
        fmt.Println("Fail to read config.ini")
        panic(err)
    }

    section := cfg.Section("")

    domains := section.Key("DOMAINS").String()
    json := jsoniter.ConfigCompatibleWithStandardLibrary

    err = json.Unmarshal([]byte(domains), &DOMAINS)
    if err != nil {
        fmt.Println("DOMAINS配置解析失败")
        panic(err)
    }

    consoleLog, err := section.Key("CONSOLE_LOG").Bool()
    if err == nil {
        CONSOLE_LOG = consoleLog
    }
    utils.CONSOLE_LOG = CONSOLE_LOG

    GRPC_HOST = section.Key("GRPC_HOST").String()

    gzipCmpr, err := section.Key("GZIP_COMPRESSION").Int()
    if err == nil {
        GZIP_COMPRESSION = gzipCmpr
    }
    coding.GZIP_COMPRESSION = GZIP_COMPRESSION

    heartBeat, err := section.Key("HEARTBEAT").Int()
    if err == nil && heartBeat > 0 {
        HEARTBEAT = heartBeat
    }

    httpTimeout, err := section.Key("HTTP_TIMEOUT").Int()
    if err == nil && httpTimeout >= 0 {
        HTTP_TIMEOUT = httpTimeout
    }
}
