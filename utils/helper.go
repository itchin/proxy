package utils

import (
    "fmt"
)

var CONSOLE_LOG bool

func ConsoleLog(str string, args ...interface{}) {
    if CONSOLE_LOG {
        if len(args) > 0 {
            fmt.Printf(str + "\n", args...)
        } else {
            fmt.Println(str)
        }
    }
}
