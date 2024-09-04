package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zachturing/login/api"
	"github.com/zachturing/login/common/define"
	"github.com/zachturing/util/config"
	"github.com/zachturing/util/config/business/common"
	"github.com/zachturing/util/database/mysql"
	"github.com/zachturing/util/database/redis"
	"github.com/zachturing/util/log"
)

func initService() error {
	err := config.Register(config.Common, common.MapEnvToConfig, define.Env)
	if err != nil {
		log.Errorf("register config failed, err:%v", err)
		return err
	}

	if err = initDatabase(); err != nil {
		return err
	}
	return initRedis()
}

func initRoute() *gin.Engine {
	router := gin.Default()
	group := router.Group("/api")
	{
		group.POST("/auth/sms", api.SendSMS)
		// group.POST("/login/phone", api.LoginPhone)
	}

	return router
}

func initDatabase() error {
	dbConfig, err := common.GetMysqlConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(dbConfig)
	return mysql.InitDatabase(dbConfig)
}

func initRedis() error {
	redisConfig, err := common.GetRedisConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(redisConfig)
	return redis.InitRedis(redisConfig)
}
