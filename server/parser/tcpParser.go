package parser

import (
    "net"

    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
    "github.com/itchin/proxy/utils/constant"
    "github.com/itchin/proxy/utils/model"
    "github.com/tidwall/gjson"
)

var TcpParser tcpParser

var RespChan = make(chan *model.Response)

type tcpParser struct{}

func (t *tcpParser) Receiver(conn net.Conn, buf []byte) {
    json := string(buf)
    utils.ConsoleLog("packet: %s", json)
    action := gjson.Get(json, "action").Int()
    data := []byte(gjson.Get(json, "data").String())
    switch action {
    case constant.ON_OPEN:
        open := new(model.Open)
        open.UnmarshalJSON(data)
        t.Open(conn, open)
    case constant.ON_MESSAGE:
        response := new(model.Response)
        response.UnmarshalJSON(data)
        t.Message(conn, response)
    }
}

// 建立连接时写入配置
func (*tcpParser) Open(conn net.Conn, open *model.Open) {
    Connection.Open(open.Domains, conn)
    utils.ConsoleLog("当前链接：%v", Connection.All())
}

// 处理http包
func (*tcpParser) Message(conn net.Conn, response *model.Response) {
    response.Body = string(coding.Decode([]byte(response.Body)))
    RespChan <- response
}
