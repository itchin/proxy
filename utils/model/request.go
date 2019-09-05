package model

type Request struct {
    Domain string `json:"domain"`
    Uri string `json:"uri"`
    Method string `json:"method"`
    Header map[string][]string `json:"header"`
    Body string `json:"body"`
}
