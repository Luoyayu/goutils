package myRedis

import "github.com/go-redis/redis/v7"

// Client is my costumed redis.Client
// it has some Base Methods, see more in methods.go
type Client struct {
	*redis.Client
}
