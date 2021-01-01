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

type GrpcParser struct {
    hp *httpParser
}

// 建立链接后向服务端注册域名
func (c *GrpcParser) Register(workerId int) {
    r := new(model.Register)
    r.Domains = getDomains()
    data, _ := r.MarshalJSON()
    GrpcClient.Write(workerId, constant.REGISTER, string(data))
}

func (c *GrpcParser) Listener(workerId int) {
    for {
        // 接收从服务端返回的数据流
        resp, err := GrpcClient.stream[workerId].Recv()
        if err == io.EOF {
            log.Println("EOF...")
            break
        }

        if err != nil {
            log.Println("receive error:", err)
        }

        // 处理来自服务端的消息
        c.Message(workerId, resp.Data)
    }
}

// 处理从服务端转发的http请求
func (c *GrpcParser) Message(workerId int, data string) {
    request := new(model.Request)
    utils.ConsoleLog(data)
    _ = request.UnmarshalJSON([]byte(data))
    if request.Domain == "" {
        return
    }
    if c.hp == nil {
        c.hp = NewHttpParser()
    }
    response, err := c.hp.Request(request)
    if err != nil {
        utils.ConsoleLog("err: %v", response, err)
        return
    }
    c.sendResponse(workerId, response)
}

// 将http response发送回tcp服务端
func (c *GrpcParser) sendResponse(workerId int, response *model.Response) {
    data, _ := response.MarshalJSON()
    GrpcClient.Write(workerId, constant.HTTP_PACKET, string(data))
}


func (*GrpcParser) Beat(workerId int) {
    heartbeat := config.HEARTBEAT
    if heartbeat > 0 {
        ticker := time.NewTicker(time.Second * time.Duration(heartbeat))
        for {
            select {
            case <- ticker.C:
                GrpcClient.Write(workerId, constant.BEAT, "")
            }
        }
    }
}
