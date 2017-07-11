package storage

import (
	"bytes"
	"encoding/gob"
)

// dataToBytes onverts a value into a slice of bytes
func dataToBytes(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// bytesToData decodes a slice of bytes into a value. data should be a pointer to the a value of the type to be decoded.
func bytesToData(data interface{}, b []byte) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(data)
}
