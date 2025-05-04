package protocol

import (
	"strings"
)

type Request struct {
	Path string
	Body []byte
}

func (m *Request) Encode() []byte {
	return []byte(m.Path + "\n" + string(m.Body))
}

func (m *Request) Decode(data []byte) {
	parts := strings.Split(string(data), "\n")
	m.Path = parts[0]
	m.Body = []byte(parts[1])
}
