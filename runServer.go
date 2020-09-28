package main

import (
    "github.com/itchin/proxy/proto"
    "github.com/itchin/proxy/server/process"
    "google.golang.org/grpc"
    "net"
    "net/http"

    "github.com/itchin/proxy/server/config"
    "github.com/itchin/proxy/server/handle"
)

func main() {
    http.HandleFunc("/", handle.HttpHandle.Router)
    go http.ListenAndServe(config.HTTP_HOST, nil)

    server := grpc.NewServer()

    proto.RegisterGrpcServer(server, &process.Streamer{})

    listener, err := net.Listen("tcp", config.TCP_HOST)
    if err != nil {
        panic(err)
    }
    if err := server.Serve(listener); err != nil {
        panic(err)
    }
}
