package main

import (
    "github.com/itchin/proxy/proto"
    "github.com/itchin/proxy/server/config"
    "github.com/itchin/proxy/server/handle"
    "github.com/itchin/proxy/server/process"
    "golang.org/x/net/netutil"
    "google.golang.org/grpc"
    "log"
    "net"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handle.HttpHandle.Router)
    s := &http.Server{
        Handler: mux,
    }
    l, err := net.Listen("tcp4", config.HTTP_HOST)
    if err != nil {
        log.Println("Http server error:", err)
    }
    go s.Serve(netutil.LimitListener(l, config.MAX_CONN))

    //http.HandleFunc("/", handle.HttpHandle.Router)
    //go http.ListenAndServe(config.HTTP_HOST, nil)

    log.Println("http server start...")

    server := grpc.NewServer()

    proto.RegisterGrpcServer(server, &process.Streamer{})
    log.Println("grpc server start...")
    listener, err := net.Listen("tcp4", config.GRPC_HOST)
    if err != nil {
        panic(err)
    }
    if err := server.Serve(listener); err != nil {
        panic(err)
    }
}
