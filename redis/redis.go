package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/newdee/aipaper-util/database"
)

var (
	RedisDB *redis.Client
)

func GetGlobalClient() *redis.Client {
	return RedisDB
}

func InitRedis(dbConfig database.DatabaseConfig) error {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", dbConfig.GetHost(), dbConfig.GetPort()),
		Password: dbConfig.GetPassword(),
		DB:       1, // login服务指定使用Redis的1号数据库
	})
	return nil
}
