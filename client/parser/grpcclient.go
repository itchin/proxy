package parser

import (
    "github.com/itchin/proxy/proto"
    "log"
    "sync"
)

var GrpcClient grpcClient

type grpcClient struct{
    mu sync.Mutex
    stream map[int]proto.Grpc_ProcessClient
}

func init()  {
    GrpcClient.stream = make(map[int]proto.Grpc_ProcessClient)
}

func (c *grpcClient) Set(i int, stream proto.Grpc_ProcessClient) {
    c.mu.Lock()
    c.stream[i] = stream
    log.Println("worker:", i)
    c.mu.Unlock()
}

func (c *grpcClient) Write(i int, data_type int32, data string) {
    err := c.stream[i].Send(&proto.Request{Type: data_type, Data: data})
    if err != nil {
        log.Println("send msg fail:", err)
    }
}
