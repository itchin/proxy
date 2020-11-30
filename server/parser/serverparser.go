package parser

import (
    "bytes"
    "github.com/itchin/proxy/proto"
    "github.com/itchin/proxy/utils/model"
    "io"
    "log"
    "net/http"
)

var ServerParser serverParser

type serverParser struct {
    //mu sync.Mutex
}

// 将请求头转发到内网服务器
func (p *serverParser) Request(httpId string,stream proto.Grpc_ProcessServer,  domain string, request *http.Request) {
    req := model.Request{
        HttpId: httpId,
        Domain: domain,
        Uri: request.RequestURI,
        Method: request.Method,
        Header: request.Header,
        Body: p.getBody(request.Body),
    }
    r, _ := req.MarshalJSON()

    //p.mu.Lock()
    err := stream.Send(&proto.Response{Data: string(r)})
    //p.mu.Unlock()
    if err != nil {
        log.Println("grpc error, send msg fail:", err)
    }
}

func (*serverParser) getBody(body io.ReadCloser) string {
    if body == nil {
        return ""
    }
    buf := new(bytes.Buffer)
    _, _ = buf.ReadFrom(body)
    return buf.String()
}
