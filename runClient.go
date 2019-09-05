package main

import (
    "fmt"
    "net"

    "github.com/itchin/proxy/client/config"
    "github.com/itchin/proxy/client/parser"
)

func main() {
    conn, err := net.Dial("tcp", config.TCP_HOST)
    if err != nil {
        fmt.Printf("connect failed, err : %v\n", err.Error())
        return
    }
    parser.Connection.Set(conn)
    defer conn.Close()

    parser.TcpParser.Open()

    go parser.TcpParser.Beat()

    parser.TcpParser.Listener()
}
