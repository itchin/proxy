package main

import (
    "context"
    "github.com/itchin/proxy/client/config"
    "github.com/itchin/proxy/client/parser"
    "github.com/itchin/proxy/proto"
    "google.golang.org/grpc"
    "log"
)

func main() {
    for i := 0; i < config.WORKERS; i ++ {
        go func(workerId int) {
            conn, err := grpc.Dial(config.GRPC_HOST, grpc.WithInsecure())
            if err != nil {
                log.Printf("connect failed, err : %v\n", err.Error())
                return
            }
            defer conn.Close()

            // 声明客户端
            client := proto.NewGrpcClient(conn)
            // 声明上下文
            ctx := context.Background()
            stream, err := client.Process(ctx)
            parser.GrpcClient.Set(workerId, stream)

            var clientParser parser.GrpcParser

            clientParser.Register(workerId)

            go clientParser.Beat(workerId)

            clientParser.Listener(workerId)
        }(i)
    }
    select {}
}
