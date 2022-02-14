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

type serverParser struct {}

// 将请求头转发到内网服务器
func (p *serverParser) Request(httpId int, stream proto.Grpc_ProcessServer, domain string, request *http.Request) (err error) {
    req := model.Request{
        HttpId: httpId,
        Domain: domain,
        Uri: request.RequestURI,
        Method: request.Method,
        Header: request.Header,
        Body: p.getBody(request.Body),
    }
    r, err := req.MarshalJSON()
    if err != nil {
        log.Println("marshl json error：", err)
        return
    }

    err = stream.Send(&proto.Response{Data: string(r)})
    if err != nil {
        log.Println("grpc error, send msg fail:", err)
    }
    return
}

//将http body的内容转为字符串格式
func (*serverParser) getBody(body io.ReadCloser) string {
    if body == nil {
        return ""
    }
    buf := new(bytes.Buffer)
    _, err := buf.ReadFrom(body)
    if err != nil {
        log.Println(err)
        return ""
    }
    return buf.String()
}
