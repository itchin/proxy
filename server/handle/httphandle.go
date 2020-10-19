package handle

import (
    "github.com/itchin/proxy/server/process"
    "github.com/itchin/proxy/utils/coding"
    "net/http"

    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils/model"
)

var HttpHandle httpHandle

type httpHandle struct{}

func (h *httpHandle) Router(response http.ResponseWriter, request *http.Request) {
    domain := parser.Addr(request.Host)
    stream := parser.Streams.Get(domain)
    if stream == nil {
        process.RespChan <- &model.Response{
            Body: "页面不存在",
        }
        return
    }
    parser.ServerParser.Request(stream, domain, request)
    h.packet(&response, <- process.RespChan, request)
}

func (*httpHandle) packet(response *http.ResponseWriter, remoteResp *model.Response, request *http.Request) {
    header := (*response).Header()
    for k, v := range remoteResp.Header {
        header.Set(k, v[0])
    }
    header.Set("Host", request.Host)
    buf := coding.Decode([]byte(remoteResp.Body))
    (*response).WriteHeader(remoteResp.StatusCode)
    _, _ = (*response).Write(buf)
}
