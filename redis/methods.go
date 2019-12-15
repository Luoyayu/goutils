package myRedis

import "github.com/go-redis/redis/v7"

func InitRedis(op *redis.Options) *Client {
	c := &Client{}
	if c.Client = redis.NewClient(op); c.Client == nil {
		panic("init redis error!")
	}
	return c
}

func (r *Client) Exists(key string) (bool, error) {
	cmd := r.Client.Exists(key)
	return cmd.Val() == 1, cmd.Err()
}
