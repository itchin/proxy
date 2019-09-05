package parser

import (
    "io"
    "time"

    "github.com/itchin/proxy/client/config"
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
    "github.com/itchin/proxy/utils/constant"
    "github.com/itchin/proxy/utils/model"
    jsoniter "github.com/json-iterator/go"
)

var TcpParser tcpParser

type tcpParser struct{}

// 建立链接后向服务端注册域名
func (*tcpParser) Open() {
    open := new(model.Open)
    open.Domains = getDomains()
    packet := new(model.Packet)
    packet.Action = constant.ON_OPEN
    packet.Data = open
    json := jsoniter.ConfigCompatibleWithStandardLibrary
    buf, _ := json.Marshal(packet)
    utils.ConsoleLog(string(buf))
    Connection.Write(buf)
}

func (t *tcpParser) Listener() {
    cache := make([]byte, 0)
    packet := make([]byte, 0)
    buf := make([]byte, 1024)
    for {
        n, err := Connection.Get().Read(buf)
        if err != nil || err == io.EOF {
            break
        }
        utils.ConsoleLog(string(buf[:n]))
        packet, cache = coding.Unpack(append(cache, buf[:n]...))

        if len(packet) > 0 {
            t.Message(packet)
        }
    }
}

// 处理从服务端转发的http请求
func (t *tcpParser) Message(buf []byte) {
    request := new(model.Request)
    request.UnmarshalJSON(buf)
    if request.Domain == "" {
        return
    }
    response, err := HttpParser.Request(request)
    if err != nil {
        utils.ConsoleLog("err: %v", response, err)
        // TODO:tcp返回错误信息
        return
    }
    t.sendResponse(response)
}

// 将http response发送回tcp服务端
func (*tcpParser) sendResponse(response *model.Response) {
    packet := new(model.Packet)
    packet.Action = constant.ON_MESSAGE
    packet.Data = response
    json := jsoniter.ConfigCompatibleWithStandardLibrary
    buf, _ := json.Marshal(packet)
    Connection.Write(buf)
}


func (*tcpParser) Beat() {
    heartbeat := config.HEARTBEAT
    if heartbeat > 0 {
        ticker := time.NewTicker(time.Second * time.Duration(heartbeat))
        for {
            select {
            case <- ticker.C:
                Connection.Get().Write([]byte("0"))
            }
        }
    }
}
