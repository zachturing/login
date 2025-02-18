package main

import (
	"fmt"
	"github.com/zachturing/login/redis"

	"github.com/gin-gonic/gin"
	"github.com/newdee/aipaper-util/config"
	"github.com/newdee/aipaper-util/config/business/common"
	"github.com/newdee/aipaper-util/database/mysql"
	"github.com/newdee/aipaper-util/log"
	"github.com/zachturing/login/api"
	"github.com/zachturing/login/common/define"
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
		group.POST("/auth/login/phone", api.LoginPhone)
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
