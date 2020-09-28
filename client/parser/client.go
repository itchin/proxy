package parser

import (
    "github.com/itchin/proxy/proto"
)

var Client client

type client struct{
    stream proto.Grpc_ProcessClient
}

func (c *client) Set(stream proto.Grpc_ProcessClient) {
    c.stream = stream
}

func (c *client) Get() proto.Grpc_ProcessClient {
    return c.stream
}

func (c *client) Write(data_type int32, data string) {
    c.stream.Send(&proto.Request{Type: data_type, Data: data})
}
