package model

type Response struct {
    Header map[string][]string `json:"header"`
    Body string `json:"body"`
}
