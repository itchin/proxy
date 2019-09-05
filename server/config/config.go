package config

import (
    "fmt"
    "os"

    "github.com/itchin/proxy/utils"
    "gopkg.in/ini.v1"
)

var (
    // 是否在控制台中打印日志
    CONSOLE_LOG = true

    TCP_HOST string

    HTTP_HOST string
)

func init() {
    cfg, err := ini.Load("server.ini")
    if err != nil {
        fmt.Printf("Fail to read config.ini: %v", err)
        os.Exit(1)
    }

    section := cfg.Section("")

    consoleLog, err := section.Key("CONSOLE_LOG").Bool()
    if err == nil {
        CONSOLE_LOG = consoleLog
    }
    utils.CONSOLE_LOG = CONSOLE_LOG

    TCP_HOST = section.Key("TCP_HOST").String()

    HTTP_HOST = section.Key("HTTP_HOST").String()
}
