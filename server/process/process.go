package process

import (
    "github.com/itchin/proxy/proto"
    "github.com/itchin/proxy/server/handle"
    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/constant"
    "github.com/itchin/proxy/utils/model"
    "io"
    "log"
    "sync"
)

type Streamer struct{}

var mu = &sync.RWMutex{}

func (s *Streamer) Process(stream proto.Grpc_ProcessServer) error {
    for {
        //调用远程代码执行后返回值
        req, err := stream.Recv()
        if err == io.EOF {
            log.Println("EOF")
            parser.Streams.Close(stream)
            log.Println("当前链接数:", parser.Streams.All())
            return nil
        }
        if err != nil {
            log.Println("read from connect failed, err:", err)
            parser.Streams.Close(stream)
            //log.Println("当前链接数:", parser.Streams.All())
            return err
        }

        //根据定义的返回值类型匹配逻辑
        switch req.Type {
        //返回http响应
        case constant.HTTP_PACKET:
            utils.ConsoleLog("receive: %v", req)
            response := new(model.Response)
            _ = response.UnmarshalJSON([]byte(req.Data))
            mu.RLock()
            c := handle.HttpHandle.Chans[response.HttpId]
            mu.RUnlock()
            c <- response

        // 注册域名与stream对象的映射
        case constant.REGISTER:
            rg := new(model.Register)
            err := rg.UnmarshalJSON([]byte(req.Data))
            if err != nil {
                log.Println("register error:", err)
                continue
            }
            parser.Streams.Register(rg.Domains, stream)
            log.Println("当前链接：", parser.Streams.All())

        // 心跳
        case constant.BEAT:

        }
    }
}
