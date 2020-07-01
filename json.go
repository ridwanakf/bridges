package bridges

// Json module used for encoding and decoding JSON string to Golang's struct
//go:generate mockgen -destination=json/jsoniterator_mock.go -package=json github.com/ridwanakf/go-bridges Json
type Json interface {
	// Unmarshal decodes JSON string to a struct
	// v should be a pointer
	Unmarshal(data []byte, v interface{}) error

	// Marshal encodes struct into JSON
	// v should not be a pointer
	Marshal(v interface{}) ([]byte, error)
}
