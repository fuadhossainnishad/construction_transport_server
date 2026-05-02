package redis

type RateLimiter struct {
	Key   string
	Value int
}
