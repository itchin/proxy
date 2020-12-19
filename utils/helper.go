package utils

import (
    "log"
    "strings"
)

var CONSOLE_LOG bool

func ConsoleLog(str string, args ...interface{}) {
    if CONSOLE_LOG {
        if len(args) > 0 {
            log.Printf(str + "\n", args...)
        } else {
            log.Println(str)
        }
    }
}

func Addr(host string) string {
    splits := strings.Split(host, ":")
    if len(splits) > 1 {
        return splits[0]
    }
    return host
}
