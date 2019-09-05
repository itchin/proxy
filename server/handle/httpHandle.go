package handle

import (
    "net/http"

    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils/model"
)

var HttpHandle httpHandle

type httpHandle struct{}

func (h *httpHandle) Router(response http.ResponseWriter, request *http.Request) {
    go func() {
        domain := parser.Addr(request.Host)
        conn := parser.Connection.Get(domain)
        if conn == nil {
            parser.RespChan <- &model.Response{
                Body: "页面不存在",
            }
            return
        }
        parser.HttpParser.Request(conn, domain, request)
    }()
    h.packet(&response, <- parser.RespChan, request.Host)
}

func (*httpHandle) packet(response *http.ResponseWriter, remoteResp *model.Response, host string) {
    header := (*response).Header()
    for k, v := range remoteResp.Header {
        header.Set(k, v[0])
    }
    header.Set("Host", host)
    (*response).Write([]byte(remoteResp.Body))
}
