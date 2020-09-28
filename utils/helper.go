package utils

import (
    "log"
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
