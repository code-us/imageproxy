package imageproxy

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"gopkg.in/redis.v5"
)

var client *redis.Client

func InitRedis() {

	uri, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	var host string
	var password string
	var db int

	switch uri.Scheme {
	case "redis":
		host = uri.Host
		if uri.User != nil {
			password, _ = uri.User.Password()
		}
		if len(uri.Path) > 1 {
			var err error
			db, err = strconv.Atoi(uri.Path[1:])
			if err != nil {
				panic("Database must be an integer")
			}
		}
	case "unix":
		host = uri.Path
	default:
		panic("invalid Redis database URI scheme")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Redis")
}

func RedisHGet(key, field string) *string {
	cmd := client.HGet(key, field)
	if cmd.Err() == redis.Nil {
		return nil
	} else if cmd.Err() != nil {
		panic(cmd.Err())
	} else {
		val := cmd.Val()
		return &val
	}
}

func RedisHDel(key string, fields ...string) {
	err := client.HDel(key, fields...).Err()
	if err != nil {
		panic(err)
	}
}

func RedisHSet(key, field, value string) {
	err := client.HSet(key, field, value).Err()
	if err != nil {
		panic(err)
	}
}

func RedisSAdd(key string, members ...interface{}) {
	err := client.SAdd(key, members...).Err()
	if err != nil {
		panic(err)
	}
}

func RedisSRem(key string, members ...interface{}) {
	err := client.SRem(key, members...).Err()
	if err != nil {
		panic(err)
	}
}

func RedisSet(key string, value interface{}, expiration time.Duration) {
	err := client.Set(key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGet(key string) *string {
	cmd := client.Get(key)
	if cmd.Err() == redis.Nil {
		return nil
	} else if cmd.Err() != nil {
		panic(cmd.Err())
	} else {
		val := cmd.Val()
		return &val
	}
}

func RedisDel(keys ...string) {
	err := client.Del(keys...).Err()
	if err != nil {
		panic(err)
	}
}

func RedisSMembers(key string) []string {
	sMembers := client.SMembers(key)

	err := sMembers.Err()
	if err != nil {
		panic(err)
	}

	return sMembers.Val()
}
