package process

import (
    "io"
    "net"

    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
)

func Accept(server net.Listener) {
    for {
        conn, err := server.Accept()
        if err != nil {
            continue
        }
        conn.Write(coding.Packet([]byte("Connection Success")))
        defer conn.Close()
        go run(conn)
    }
}

func run(conn net.Conn) {
    cache := make([]byte, 0)// 缓存被截断的数据流
    packet := make([]byte, 0)// 完整的数据包
    buf := make([]byte, 1024)// 单次接收的数据流

    for {
        n, err := conn.Read(buf)

        if err != nil || err == io.EOF {
            parser.Connection.Close(conn)
            utils.ConsoleLog("read from connect failed, err: %v", err)
            utils.ConsoleLog("当前链接：%v", parser.Connection.All())
            break
        }

        if string(buf[:n]) == "0" {
            continue
        }

        utils.ConsoleLog("receive: %s", string(buf[:n]))
        packet, cache = coding.Unpack(append(cache, buf[:n]...))
        if len(packet) > 0 {
            parser.TcpParser.Receiver(conn, packet)
        }
    }
}
