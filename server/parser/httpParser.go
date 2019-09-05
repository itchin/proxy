package parser

import (
    "bytes"
    "io"
    "net"
    "net/http"

    "github.com/itchin/proxy/utils/coding"
    "github.com/itchin/proxy/utils/model"
)

var HttpParser httpParser

type httpParser struct{}

func (p *httpParser) Request(conn net.Conn,domain string, request *http.Request) {
    req := model.Request{
        Domain: domain,
        Uri: request.RequestURI,
        Method: request.Method,
        Header: request.Header,
        Body: p.getBody(request.Body),
    }
    r, _ := req.MarshalJSON()

    conn.Write(coding.Packet(r))
}

func (*httpParser) getBody(body io.ReadCloser) string {
    if body == nil {
        return ""
    }
    buf := new(bytes.Buffer)
    buf.ReadFrom(body)
    return buf.String()
}
