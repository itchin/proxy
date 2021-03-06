package model

type Response struct {
    HttpId int `json:"http_id"`
    Gzip bool `json:"gzip_level"`
    StatusCode int `json:"status_code"`
    Header map[string][]string `json:"header"`
    Body string `json:"body"`
}
