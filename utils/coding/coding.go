package coding

import (
    "encoding/base64"
)

func Encode(buf []byte) string {
    buf, _ = GzipEncode(buf)
    return base64.StdEncoding.EncodeToString(buf)
}

func Decode(code []byte, gzip bool) (buf []byte, err error) {
    buf, _ = base64.StdEncoding.DecodeString(string(code))
    if gzip {
        buf, err = GzipDecode(buf)
    }
    return
}
