package bridges

// RedisJson adds the capability to Redis by adding capabilities to Unmarshal and Marshal JSON on this very module
//go:generate mockgen -destination=redisjson/redisjson_mock.go -package=redisjson github.com/ridwanakf/go-bridges RedisJson
type RedisJson interface {
	Redis
	// GetUnmarshalled fetches value from redis key
	// v should be pointer
	GetUnmarshalled(key string, v interface{}) error

	// SetexMarshalled sets value with expiry
	// v should not be a pointer
	SetexMarshalled(key string, seconds int, v interface{}) error
}
