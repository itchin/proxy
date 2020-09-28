package main

import (
    "context"
    "fmt"
    "github.com/itchin/proxy/client/config"
    "github.com/itchin/proxy/client/parser"
    "github.com/itchin/proxy/proto"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial(config.TCP_HOST, grpc.WithInsecure())
    if err != nil {
        fmt.Printf("connect failed, err : %v\n", err.Error())
        return
    }
    defer conn.Close()

    // 声明客户端
    client := proto.NewGrpcClient(conn)
    // 声明上下文
    ctx := context.Background()
    stream, err := client.Process(ctx)
    parser.Client.Set(stream)

    parser.ClientParser.Register()

    go parser.ClientParser.Beat()

    parser.ClientParser.Listener()
}
