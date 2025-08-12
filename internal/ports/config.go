package ports

import "time"

type Config interface {
	Sub(key string) Config
	GetString(key string) string
	GetStringSlice(key string) []string
	GetInt(key string) int
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	SetDefault(key string, value interface{})
}
