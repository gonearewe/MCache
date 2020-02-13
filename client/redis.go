package client

import (
	"github.com/go-redis/redis"
)

// for redis client, we take advantage of an existed repo,
// all we need is to wrap it.

type redisClient struct {
	*redis.Client
}

func (r *redisClient) get(key string) (string, error) {
	res, e := r.Get(key).Result()
	if e == redis.Nil {
		return "", nil
	}
	return res, e
}

func (r *redisClient) set(key, value string) error {
	return r.Set(key, value, 0).Err()
}

func (r *redisClient) del(key string) error {
	return r.Del(key).Err()
}

func (r *redisClient) Run(req *Request) {
	switch req.Type {
	case RequestGet:
		var v string
		v, req.Error = r.get(req.Key)
		req.Val = []byte(v)
		return

	case RequestSet:
		req.Error = r.set(req.Key, string(req.Val))
		return

	case RequestDel:
		req.Error = r.del(req.Key)
		return
	}

	panic("unknown request name " + req.Type)
}

func (r *redisClient) PipelinedRun(cmds []*Request) {
	if len(cmds) == 0 {
		return
	}
	pipe := r.Pipeline()
	cmders := make([]redis.Cmder, len(cmds))
	for i, c := range cmds {
		if c.Type == RequestGet {
			cmders[i] = pipe.Get(c.Key)
		} else if c.Type == RequestSet {
			cmders[i] = pipe.Set(c.Key, c.Val, 0)
		} else if c.Type == RequestDel {
			cmders[i] = pipe.Del(c.Key)
		} else {
			panic("unknown cmd name " + c.Type)
		}
	}

	_, err := pipe.Exec()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	// verify Get type
	for i, c := range cmds {
		if c.Type == RequestGet {
			value, e := cmders[i].(*redis.StringCmd).Result()
			if e == redis.Nil {
				value, e = "", nil
			}
			c.Val, c.Error = []byte(value), e
		} else {
			c.Error = cmders[i].Err()
		}
	}
}

func newRedisClient(server string) *redisClient {
	return &redisClient{redis.NewClient(&redis.Options{Addr: server + ":6379", ReadTimeout: -1})}
}
