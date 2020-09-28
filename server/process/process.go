package process

import (
    "github.com/itchin/proxy/proto"
    "github.com/itchin/proxy/utils/model"
    "io"

    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/constant"
)

var RespChan = make(chan *model.Response)

type Streamer struct{}

func (s *Streamer) Process(stream proto.Grpc_ProcessServer) error {
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            parser.Streams.Close(stream)
            utils.ConsoleLog("EOF")
            utils.ConsoleLog("当前链接：%v", parser.Streams.All())
            return nil
        }
        if err != nil {
            parser.Streams.Close(stream)
            utils.ConsoleLog("read from connect failed, err: %v", err)
            utils.ConsoleLog("当前链接：%v", parser.Streams.All())
            return err
        }

        switch req.Type {

        case constant.HTTP_PACKET:
            utils.ConsoleLog("receive: %v", req)
            response := new(model.Response)
            _ = response.UnmarshalJSON([]byte(req.Data))
            RespChan <- response

        // 注册域名与stream对象的映射
        case constant.REGISTER:
            rg := new(model.Register)
            err := rg.UnmarshalJSON([]byte(req.Data))
            if err != nil {
                utils.ConsoleLog("register error: %v", err)
                continue
            }
            parser.Streams.Register(rg.Domains, stream)
            utils.ConsoleLog("当前链接：%v", parser.Streams.All())

        // 心跳
        case constant.BEAT:

        }
    }
}
