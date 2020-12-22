package handle

import (
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
    cmap "github.com/orcaman/concurrent-map"
    "log"
    "net/http"
    "strconv"
    "sync"

    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils/model"
)

var HttpHandle httpHandle

type httpHandle struct {
    // response映射表
    CMap cmap.ConcurrentMap
    mu sync.Mutex
    httpId int
}

func init() {
    HttpHandle.CMap = cmap.New()
}

func (h *httpHandle) Router(rw http.ResponseWriter, request *http.Request) {
    h.mu.Lock()
    seq := Capacity.Shift()
    httpId := strconv.Itoa(h.httpId)
    h.CMap.Set(httpId, make(chan *model.Response))
    h.httpId++
    h.mu.Unlock()

    domain := utils.Addr(request.Host)
    stream := parser.Streams.Get(domain)
    if stream == nil {
        c, _ := h.CMap.Get(httpId)
        c.(chan *model.Response) <- &model.Response{
            Body: "页面不存在",
        }
        return
    } else {
        parser.ServerParser.Request(httpId, stream, domain, request)
    }

    remoteResp := h.listener(httpId)
    h.responseHandle(remoteResp, &rw, request, seq)
}

func (h *httpHandle) listener(httpId string) (remoteResp *model.Response) {
    ci, _ := h.CMap.Get(httpId)
    c := ci.(chan *model.Response)
    remoteResp = <- c
    close(c)
    h.CMap.Remove(httpId)
    return
}

func (h *httpHandle) responseHandle(remoteResp *model.Response, rw *http.ResponseWriter, request *http.Request, seq int) {
    for k, v := range remoteResp.Header {
        (*rw).Header().Set(k, v[0])
    }
    (*rw).Header().Set("Host", request.Host)
    buf, err := coding.Decode([]byte(remoteResp.Body))
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
