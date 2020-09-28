package parser

import (
    "github.com/itchin/proxy/proto"
    "strings"
    "sync"

    "github.com/itchin/proxy/utils"
)

var Streams streams

type streams struct{
    // domain => stream
    m map[string]proto.Grpc_ProcessServer
    mu sync.RWMutex
}

// 获取域名对应的链接
func (s *streams) Get(domain string) proto.Grpc_ProcessServer {
    s.mu.Lock()
    if stream, ok := s.m[domain]; ok {
        s.mu.Unlock()
        return stream
    }
    return nil
}

// 建立连接后，将域名与grpc流对象绑定
func (s *streams) Register(domains []string, stream proto.Grpc_ProcessServer) {
    s.mu.Lock()
    if str, isExists := s.isExists(domains); isExists {
        s.mu.Unlock()
        utils.ConsoleLog(str)
        return
    }

    if s.m == nil {
        s.m = make(map[string]proto.Grpc_ProcessServer)
    }
    for _, domain := range domains {
        s.m[domain] = stream
    }
    s.mu.Unlock()
}

// 判断客户端请求注册的host，是否已存在于服务端注册信息中
func (c *streams) isExists(domains []string) (string, bool) {
    dms := make([]string, 0)
    for _, domain := range domains {
        dm := domain
        if _, ok := c.m[dm]; ok {
            dms = append(dms, dm)
        }
    }
    if len(dms) > 0 {
        str := "操作失败，以下Host已在服务器中注册：" + strings.Join(dms, ",")
        return str, true
    }
    return "", false
}

// 客户端断开连接时，将其从链接池移除
func (c *streams) Close(stream proto.Grpc_ProcessServer) {
    c.mu.Lock()
    for k, v := range c.m {
        if v == stream {
            c.mu.Lock()
            delete(c.m, k)
            c.mu.Unlock()
        }
    }
}

func (c *streams) All() map[string]proto.Grpc_ProcessServer {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.m
}
