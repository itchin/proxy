package parser

import (
    "bytes"
    "github.com/itchin/proxy/proto"
    "io"
    "net/http"

    "github.com/itchin/proxy/utils/model"
)

var HttpParser httpParser

type httpParser struct{}

func (p *httpParser) Request(stream proto.Grpc_ProcessServer,domain string, request *http.Request) {
    req := model.Request{
        Domain: domain,
        Uri: request.RequestURI,
        Method: request.Method,
        Header: request.Header,
        Body: p.getBody(request.Body),
    }
    r, _ := req.MarshalJSON()

    _ = stream.Send(&proto.Response{Data: string(r)})
}

func (*httpParser) getBody(body io.ReadCloser) string {
    if body == nil {
        return ""
    }
    buf := new(bytes.Buffer)
    _, _ = buf.ReadFrom(body)
    return buf.String()
}
