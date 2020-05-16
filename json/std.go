package json

import "encoding/json"

type std struct {
}

// NewStd construct a Golang-standard JSON library
func NewStd() *std {
	return &std{}
}

func (*std) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (*std) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
