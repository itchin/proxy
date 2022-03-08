package handle

import (
    "context"
    "github.com/itchin/proxy/server/config"
    "github.com/itchin/proxy/server/parser"
    "github.com/itchin/proxy/utils"
    "github.com/itchin/proxy/utils/coding"
    "github.com/itchin/proxy/utils/model"
    "log"
    "net/http"
    "sync"
    "time"
)

var HttpHandle httpHandle

type httpHandle struct {
    Chans []chan *model.Response
    mu *sync.RWMutex
}

func init() {
    HttpHandle.mu = &sync.RWMutex{}
    //使用channel控制活跃协程数量
    HttpHandle.Chans = make([]chan *model.Response, config.MAX_ACTIVE)
    for i := 0; i < config.MAX_ACTIVE; i++ {
        HttpHandle.Chans[i] = make(chan *model.Response)
    }
}

//路由器入口
func (h *httpHandle) Router(rw http.ResponseWriter, request *http.Request) {
    seq := Capacity.Shift()

    domain := utils.Addr(request.Host)
    stream := parser.Streams.Get(domain)
    if stream == nil {
        //如果域名为空指向，返回404
        h.mu.RLock()
        c := h.Chans[seq]
        h.mu.RUnlock()
        c <- &model.Response {
            StatusCode: 404,
            Body: "页面不存在",
        }
        Capacity.Push(seq)
        return
    } else {
        err := parser.ServerParser.Request(seq, stream, domain, request)
        if err != nil {
            log.Println("请求失败：", err)
        }
    }

    remoteResp := h.listener(seq)
    if remoteResp.Header == nil {
        h.mu.RLock()
        c := h.Chans[seq]
        h.mu.RUnlock()
        c <- remoteResp
        Capacity.Push(seq)
        return
    }
    h.responseHandle(remoteResp, &rw, request, seq)
}

//监听远程http响应
func (h *httpHandle) listener(httpId int) (remoteResp *model.Response) {
    c := h.Chans[httpId]
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.HTTP_TIMEOUT) * time.Second)
    select {
    case <- ctx.Done():
        remoteResp = &model.Response{
            StatusCode: 408,
            Body: "请求超时",
        }
    default:
        remoteResp = <- c
        cancel()
    }
    return
}

//解析响应的http，替换header
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
