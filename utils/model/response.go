package model

type Response struct {
    StatusCode int `json:"status_code"`
    Header map[string][]string `json:"header"`
    Body string `json:"body"`
}
