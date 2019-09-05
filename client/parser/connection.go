package parser

import (
    "net"

    "github.com/itchin/proxy/utils/coding"
)

var Connection connection

type connection struct{
    conn net.Conn
}

func (c *connection) Set(conn net.Conn) {
    c.conn = conn
}

func (c *connection) Get() net.Conn {
    return c.conn
}

func (c *connection) Write(buf []byte) {
    c.conn.Write(coding.Packet(buf))
}
