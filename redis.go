package RedisDataStore

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis"
	faasflow "github.com/s8sg/faas-flow"
)

type RedisDataStore struct {
	bucketName  string
	redisClient redis.UniversalClient
}

// InitFromEnv Initialize a minio DataStore object based on configuration
// Depends on s3_url, s3-secret-key, s3-access-key, s3_region(optional), workflow_name
func InitFromEnv() (faasflow.DataStore, error) {
	ds := &RedisDataStore{}

	endpoint := os.Getenv("redis_url")
	addrs := strings.Split(endpoint, ",")
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName: os.Getenv("redis_master"),
		Addrs:      addrs,
	})
	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}

	ds.redisClient = client
	return ds, nil
}

func (this *RedisDataStore) Configure(flowName string, requestId string) {
	bucketName := fmt.Sprintf("faasflow-%s-%s", flowName, requestId)

	this.bucketName = bucketName
}

func (this *RedisDataStore) Init() error {
	if this.redisClient == nil {
		return fmt.Errorf("redis client not initialized, use GetRedisDataStore()")
	}

	return nil
}

func (this *RedisDataStore) Set(key string, value string) error {
	if this.redisClient == nil {
		return fmt.Errorf("redis client not initialized, use GetRedisDataStore()")
	}

	fullPath := getPath(this.bucketName, key)
	_, err := this.redisClient.Set(fullPath, value, 0).Result()
	if err != nil {
		return fmt.Errorf("error writing: %s, error: %s", fullPath, err.Error())
	}

	return nil
}

func (this *RedisDataStore) Get(key string) (string, error) {
	if this.redisClient == nil {
		return "", fmt.Errorf("redis client not initialized, use GetRedisDataStore()")
	}

	fullPath := getPath(this.bucketName, key)
	value, err := this.redisClient.Get(fullPath).Result()
	if err != nil {
		return "", fmt.Errorf("error reading: %s, error: %s", fullPath, err.Error())
	}
	return value, nil
}

func (this *RedisDataStore) Del(key string) error {
	if this.redisClient == nil {
		return fmt.Errorf("redis client not initialized, use GetRedisDataStore()")
	}

	fullPath := getPath(this.bucketName, key)
	_, err := this.redisClient.Del(fullPath).Result()
	if err != nil {
		return fmt.Errorf("error removing: %s, error: %s", fullPath, err.Error())
	}
	return nil
}

func (this *RedisDataStore) Cleanup() error {
	key := this.bucketName + ".*"
	client := this.redisClient
	var rerr error

	iter := client.Scan(0, key, 0).Iterator()
	for iter.Next() {
		err := client.Del(iter.Val()).Err()
		if err != nil {
			rerr = err
		}
	}

	if err := iter.Err(); err != nil {
		rerr = err
	}
	return rerr
}

// getPath produces a string as bucketname.value
func getPath(bucket, key string) string {
	fileName := fmt.Sprintf("%s.value", key)
	return fmt.Sprintf("%s.%s", bucket, fileName)
}
