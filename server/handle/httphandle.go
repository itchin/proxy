package handle

import (
    "github.com/itchin/proxy/server/config"
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
    "log"
    "net/http"
    "sync"

    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils/model"
)

var HttpHandle httpHandle

type httpHandle struct {
    Chans []chan *model.Response
    mu *sync.Mutex
}

func init() {
    HttpHandle.mu = &sync.Mutex{}
    HttpHandle.Chans = make([]chan *model.Response, config.MAX_ACTIVE)
    for i := 0; i < config.MAX_ACTIVE; i++ {
        HttpHandle.Chans[i] = make(chan *model.Response)
    }
}

func (h *httpHandle) Router(rw http.ResponseWriter, request *http.Request) {
    h.mu.Lock()
    seq := Capacity.Shift()
    h.mu.Unlock()

    domain := utils.Addr(request.Host)
    stream := parser.Streams.Get(domain)
    if stream == nil {
        h.mu.Lock()
        c := h.Chans[seq]
        h.mu.Unlock()
        c <- &model.Response{
            Body: "页面不存在",
        }
        return
    } else {
        parser.ServerParser.Request(seq, stream, domain, request)
    }

    remoteResp := h.listener(seq)
    h.responseHandle(remoteResp, &rw, request, seq)
}

func (h *httpHandle) listener(httpId int) (remoteResp *model.Response) {
    c := h.Chans[httpId]
    remoteResp = <- c
    return
}

func (h *httpHandle) responseHandle(remoteResp *model.Response, rw *http.ResponseWriter, request *http.Request, seq int) {
    for k, v := range remoteResp.Header {
        (*rw).Header().Set(k, v[0])
    }
    (*rw).Header().Set("Host", request.Host)
    buf, err := coding.Decode([]byte(remoteResp.Body), remoteResp.Gzip)
    if err != nil {
        log.Println("decode error:", err)
    }
    (*rw).WriteHeader(remoteResp.StatusCode)
    _, err = (*rw).Write(buf)
    if err != nil {
        log.Println("ResponseWriter error:", err)
    }
    Capacity.Push(seq)
}
