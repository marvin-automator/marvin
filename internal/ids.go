package internal

import (
	"bytes"
	"encoding/base64"
	"github.com/gobuffalo/uuid"
)

func NewId() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(make([]byte, 0, 32))
	e := base64.NewEncoder(base64.RawURLEncoding, b)
	_, err = e.Write(u.Bytes())
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
