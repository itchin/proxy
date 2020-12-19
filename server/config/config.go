package config

import (
    "fmt"
    "os"

    "github.com/itchin/proxy/utils"
    "gopkg.in/ini.v1"
)

var (
    MAX_CONN = 1024

    MAX_ACTIVE = 20
    // 是否在控制台中打印日志
    CONSOLE_LOG = true

    GRPC_HOST string

    HTTP_HOST string
)

func init() {
    cfg, err := ini.Load("server.ini")
    if err != nil {
        fmt.Printf("Fail to read config.ini: %v", err)
        os.Exit(1)
    }

    section := cfg.Section("")

    maxConn, err := section.Key("MAX_CONN").Int()
    if err == nil {
        MAX_CONN = maxConn
    }

    maxActive, err := section.Key("MAX_ACTIVE").Int()
    if err == nil {
        MAX_ACTIVE = maxActive
    }

    consoleLog, err := section.Key("CONSOLE_LOG").Bool()
    if err == nil {
        CONSOLE_LOG = consoleLog
    }
    utils.CONSOLE_LOG = CONSOLE_LOG

    GRPC_HOST = section.Key("GRPC_HOST").String()

    HTTP_HOST = section.Key("HTTP_HOST").String()
}
