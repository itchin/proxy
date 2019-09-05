package model

type Packet struct {
    Action int8 `json:"action"`
    Data interface{} `json:"data"`
}
