package coding

import (
    "encoding/base64"
)

func Encode(buf []byte) string {
    gzipCode, _ := GzipEncode(buf)
    return base64.StdEncoding.EncodeToString(gzipCode)
}

func Decode(code []byte) (buf []byte, err error) {
    gzipCode, _ := base64.StdEncoding.DecodeString(string(code))
    buf, err = GzipDecode(gzipCode)
    return
}
