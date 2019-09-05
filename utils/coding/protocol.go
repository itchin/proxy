package coding

import (
    "bytes"
    "encoding/binary"
)

const (
    HEADER_LEN = 4
)

//封包
func Packet(buf []byte) []byte {
    return append(IntToBytes(len(buf)), buf...)
}

// 解包
func Unpack(buf []byte) ([]byte, []byte) {
    length := len(buf)

    messageLength := BytesToInt(buf[:HEADER_LEN])
    total := HEADER_LEN + messageLength
    if length < total {
        return []byte{}, buf
    } else if (length == total) {
        return buf[HEADER_LEN:], []byte{}
    } else {
        return buf[HEADER_LEN:total], buf[total:]
    }
}

//整形转换成字节
func IntToBytes(n int) []byte {
    x := int32(n)

    bytesBuffer := bytes.NewBuffer([]byte{})
    binary.Write(bytesBuffer, binary.BigEndian, x)
    return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
    bytesBuffer := bytes.NewBuffer(b)

    var x int32
    binary.Read(bytesBuffer, binary.BigEndian, &x)

    return int(x)
}
