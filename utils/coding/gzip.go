package coding

import (
    "bytes"
    "compress/gzip"
    "io/ioutil"
)

var GZIP_COMPRESSION = 5

func GzipEncode(in []byte) ([]byte, error) {
    if GZIP_COMPRESSION == 0 {
        return in, nil
    }

    var (
        buf bytes.Buffer
        out    []byte
        err    error
    )

    writer, _ := gzip.NewWriterLevel(&buf, GZIP_COMPRESSION)
    _, err = writer.Write(in)
    if err != nil {
        writer.Close()
        return out, err
    }
    err = writer.Close()
    if err != nil {
        return out, err
    }

    return buf.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
    reader, err := gzip.NewReader(bytes.NewReader(in))
    if err != nil {
        var out []byte
        return out, err
    }
    defer reader.Close()

    return ioutil.ReadAll(reader)
}
