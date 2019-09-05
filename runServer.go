package main

import (
    "net"
    "net/http"

    "github.com/itchin/proxy/server/config"
    "github.com/itchin/proxy/server/handle"
    "github.com/itchin/proxy/server/process"
)

func main() {
    http.HandleFunc("/", handle.HttpHandle.Router)
    go http.ListenAndServe(config.HTTP_HOST, nil)

    server, err := net.Listen("tcp", config.TCP_HOST)
    if err != nil {
        panic(err)
    }
    defer server.Close()

    process.Accept(server)
}
