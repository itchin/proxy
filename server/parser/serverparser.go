package parser

import (
    "bytes"
    "github.com/itchin/proxy/proto"
    "io"
    "log"
    "net/http"

    "github.com/itchin/proxy/utils/model"
)

var ServerParser serverParser

type serverParser struct{}

// 将请求头转发到内网服务器
func (p *serverParser) Request(stream proto.Grpc_ProcessServer,domain string, request *http.Request) {
    req := model.Request{
        Domain: domain,
        Uri: request.RequestURI,
        Method: request.Method,
        Header: request.Header,
        Body: p.getBody(request.Body),
    }
    r, _ := req.MarshalJSON()

    err := stream.Send(&proto.Response{Data: string(r)})
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
