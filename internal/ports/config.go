package ports

import "time"

type Config interface {
	Sub(key string) Config
	GetString(key string) string
	GetStringSlice(key string) []string
	GetInt(key string) int
	GetFloat64(key string) float64
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	SetDefault(key string, value interface{})
}
