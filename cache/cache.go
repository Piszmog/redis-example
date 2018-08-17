package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Client interface {
	Insert(id string, object interface{}) error
	Get(id string, interfaceType interface{}) error
	GetAll() (map[string]string, error)
	Remove(id string) error
	RemoveAll() error
	Close()
}

type RedisClient struct {
	client *redis.Client
	key    string
}

func CreateLocalRedisClient(key string) *RedisClient {
	return &RedisClient{client: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}), key: key}
}

func CreateRedisClient(host string, port string, password string, key string) *RedisClient {
	return &RedisClient{client: redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password, // no password set
		DB:       0,        // use default DB
	}), key: key}
}

func (redisClient RedisClient) Insert(id string, object interface{}) error {
	err := redisClient.client.HSet(redisClient.key, id, object).Err()
	if err != nil {
		errors.Wrapf(err, "failed to insert %s to %s", id, redisClient.key)
	}
	return nil
}

func (redisClient RedisClient) Get(id string, interfaceType interface{}) error {
	bytes, err := redisClient.client.HGet(redisClient.key, id).Bytes()
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve %s from %s", id, redisClient.key)
	}
	err = json.Unmarshal(bytes, interfaceType)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal %s", id)
	}
	return nil
}

func (redisClient RedisClient) GetAll() (map[string]string, error) {
	stringResults, err := redisClient.client.HGetAll(redisClient.key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve all entries in key %s", redisClient.key)
	}
	return stringResults, nil
}

func (redisClient RedisClient) Remove(id string) error {
	err := redisClient.client.HDel(redisClient.key, id).Err()
	if err != nil {
		return errors.Wrapf(err, "failed to remove %s from %s", id, redisClient.key)
	}
	return nil
}

func (redisClient RedisClient) RemoveAll() error {
	err := redisClient.client.Del(redisClient.key).Err()
	if err != nil {
		errors.Wrapf(err, "failed to remove all from %s", redisClient.key)
	}
	return nil
}

func (redisClient RedisClient) Close() {
	redisClient.client.Close()
}
