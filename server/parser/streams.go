package parser

import (
    "github.com/itchin/proxy/proto"
    "sync"
)

var Streams streams

type conn struct {
    streams []proto.Grpc_ProcessServer
    curr int
    length int
}

// 指向下一个连接
func (c *conn) next() {
    if c.curr >= c.length - 1 {
        c.curr = 0
    } else {
        c.curr++
    }
}

type streams struct{
    // domain => conn
    dc map[string]*conn
    sd map[proto.Grpc_ProcessServer]string
    mu sync.RWMutex
}

func init() {
    Streams.dc = make(map[string]*conn)
    Streams.sd = make(map[proto.Grpc_ProcessServer]string)
}

// 获取域名对应的链接
func (s *streams) Get(domain string) proto.Grpc_ProcessServer {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if conn, ok := s.dc[domain]; ok {
        // 轮询使用连接对象
        conn.next()
        stream := conn.streams[conn.curr]
        return stream
    }
    return nil
}

// 建立连接后，将域名与grpc流对象绑定
func (s *streams) Register(domains []string, stream proto.Grpc_ProcessServer) {
    s.mu.Lock()
    for _, domain := range domains {
        s.sd[stream] = domain
        // 初始化
        if _, ok := s.dc[domain]; !ok {
            s.dc[domain] = new(conn)
            s.dc[domain].streams = make([]proto.Grpc_ProcessServer, 0)
        }
        c := s.dc[domain]
        c.streams = append(c.streams, stream)
        c.length ++
    }
    s.mu.Unlock()
}

// 客户端断开连接时，将其从链接池移除
func (s *streams) Close(stream proto.Grpc_ProcessServer) string {
    s.mu.Lock()
    domain := s.sd[stream]
    delete(s.sd, stream)
    c := s.dc[domain]
    for k, v := range c.streams {
       if v == stream {
           if c.length == 1 {
               delete(s.dc, domain)
           } else if c.length == k + 1 {
               c.streams = append(c.streams[0:k - 1])
               c.length --
           } else {
               c.streams = append(c.streams[0:k], c.streams[k + 1:]...)
               c.length --
           }
           break
       }
    }
    s.mu.Unlock()
    return domain
}

func (s *streams) All() map[string]*conn {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.dc
}
