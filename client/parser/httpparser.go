package parser

import (
    "github.com/itchin/proxy/client/config"
    "github.com/itchin/proxy/utils/coding"
    "github.com/itchin/proxy/utils/model"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "time"
)

var HttpParser httpParser

type httpParser struct{}

func (h *httpParser) Request(request *model.Request) (resp *model.Response, err error) {
    var locDomain string
    if val, ok := config.DOMAINS[request.Domain]; ok {
        locDomain = val
    } else {
        return
    }
    h.headerReplace(request, locDomain)

    req, err := http.NewRequest(request.Method, locDomain + request.Uri, strings.NewReader(request.Body))
    if err != nil {
        return
    }

    for k, v := range request.Header {
        req.Header.Set(k, v[0])
    }

    var client *http.Client
    if config.HTTP_TIMEOUT > 0 {
        client = &http.Client{
            Timeout: time.Duration(30 * time.Second),
        }
    } else {
        client = &http.Client{}
    }


    response, err := client.Do(req)
    if err != nil {
        log.Println("http error:", err, "status code:", response.StatusCode, "request path:", locDomain + request.Uri)
        return
    }
    defer response.Body.Close()

    bodyByte, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println("read http response body error")
        return
    }

    resp = new(model.Response)
    resp.StatusCode = response.StatusCode
    resp.Header = response.Header
    resp.Body = coding.Encode(bodyByte)
    return
}

func (*httpParser) headerReplace(request *model.Request, locDomain string) {
    host := strings.Split(locDomain, "//")
    if referer, ok := request.Header["Referer"]; ok {
        sp := strings.Split(referer[0], "/")
        r := strings.Replace(referer[0], sp[2], host[1], 1)
        request.Header["Referer"] = []string{r}
    }
}
