package parser

import (
    "io"
    "log"
    "time"

    "github.com/itchin/proxy/client/config"
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/constant"
    "github.com/itchin/proxy/utils/model"
)

var ClientParser clientParser

type clientParser struct{}

// 建立链接后向服务端注册域名
func (*clientParser) Register() {
    r := new(model.Register)
    r.Domains = getDomains()
    data, _ := r.MarshalJSON()
    Client.Write(constant.REGISTER, string(data))
}

func (c *clientParser) Listener() {
    for {
        // 接收从服务端返回的数据流
        resp, err := Client.stream.Recv()
        if err == io.EOF {
            log.Println("EOF...")
            break
        }

        if err != nil {
            log.Println("receive error:", err)
        }

        // 处理来自服务端的消息
        c.Message(resp.Data)
    }
}

// 处理从服务端转发的http请求
func (c *clientParser) Message(data string) {
    request := new(model.Request)
    utils.ConsoleLog(data)
    request.UnmarshalJSON([]byte(data))
    if request.Domain == "" {
        return
    }
    response, err := HttpParser.Request(request)
    if err != nil {
        utils.ConsoleLog("err: %v", response, err)
        return
    }
    c.sendResponse(response)
}

// 将http response发送回tcp服务端
func (*clientParser) sendResponse(response *model.Response) {
    data, _ := response.MarshalJSON()
    Client.Write(constant.HTTP_PACKET, string(data))
}


func (*clientParser) Beat() {
    heartbeat := config.HEARTBEAT
    if heartbeat > 0 {
        ticker := time.NewTicker(time.Second * time.Duration(heartbeat))
        for {
            select {
            case <- ticker.C:
                Client.Write(constant.BEAT, "")
            }
        }
    }
}
