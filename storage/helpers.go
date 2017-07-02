package storage

import (
	"bytes"
	"encoding/gob"
)

func dataToBytes(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func bytesToData(data interface{}, b []byte) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(data)
}
