package model

type Request struct {
    HttpId int `json:"http_id"`
    Domain string `json:"domain"`
    Uri string `json:"uri"`
    Method string `json:"method"`
    Header map[string][]string `json:"header"`
    Body string `json:"body"`
}
