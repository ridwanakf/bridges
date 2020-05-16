package json

import "github.com/json-iterator/go"

type jsoniterator struct {
	json jsoniter.API
}

// NewJsoniter constructs new JSON module using Jsoniter library
func NewJsoniter() *jsoniterator {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return &jsoniterator{json: json}
}

func (j *jsoniterator) Unmarshal(data []byte, v interface{}) error {
	return j.json.Unmarshal(data, v)
}

func (j *jsoniterator) Marshal(v interface{}) ([]byte, error) {
	return j.json.Marshal(v)
}
