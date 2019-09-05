package parser

import (
    "strings"
)

func Addr(host string) string {
    splits := strings.Split(host, ":")
    if len(splits) > 1 {
        return splits[0]
    }
    return host
}