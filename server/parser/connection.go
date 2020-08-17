package parser

import (
    "net"
    "strings"
    "sync"

    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
)

var Connection connection

type connection struct{
    // domain => *conn
    conns map[string]net.Conn

    mu sync.RWMutex
}

// 获取域名对应的链接
func (c *connection) Get(domain string) net.Conn {
    c.mu.Lock()
    defer c.mu.Unlock()
    if conn, ok := c.conns[domain]; ok {
        return conn
    }
    return nil
}

// 建立连接后，与进行服务端进行初始化配置
func (c *connection) Open(domains []string, conn net.Conn) {
    c.mu.Lock()
    if str, isExists := c.isExists(domains); isExists {
        utils.ConsoleLog(str)
        conn.Write(coding.Packet([]byte(str)))
        c.mu.Unlock()
        return
    }

    if c.conns == nil {
        c.conns = make(map[string]net.Conn)
    }
    for _, domain := range domains {
        c.conns[domain] = conn
    }
    c.mu.Unlock()
}

// 判断客户端请求注册的host，是否已存在于服务端注册信息中
func (c *connection) isExists(domains []string) (string, bool) {
    dms := make([]string, 0)
    for _, domain := range domains {
        dm := domain
        if _, ok := c.conns[dm]; ok {
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
func (c *connection) Close(conn net.Conn) {
    c.mu.Lock()
    conns := make(map[string]net.Conn)
    for k, v := range c.conns {
        if v != conn {
            conns[k] = v
        }
    }
    c.conns = conns
    c.mu.Unlock()
}

func (c *connection) All() map[string]net.Conn {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.conns
}
