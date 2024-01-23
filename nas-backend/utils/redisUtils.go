package utils

import (
	"github.com/go-redis/redis"
	"log"
)

func DeleteKeysByPattern(client *redis.Client, pattern string) error {
	log.Printf("正在删除redis中以 %s 开头的key", pattern)
	keys, err := client.Keys(pattern + "*").Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := client.Del(key).Err(); err != nil {
			return err
		}
		log.Printf("已删除redis中key为 %s 的缓存", key)
	}
	return nil
}
