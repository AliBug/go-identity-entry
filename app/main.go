package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_tokenRepo "github.com/alibug/go-identity-entry/token/repository/redisdb"
	_tokenUseCase "github.com/alibug/go-identity-entry/token/usecase"
	_userHttpDelivery "github.com/alibug/go-identity-entry/user/delivery/restgin"
	_userRepo "github.com/alibug/go-identity-entry/user/repository/mongodb"
	_userUseCase "github.com/alibug/go-identity-entry/user/usecase"
	"github.com/alibug/go-identity-utils/config"
	"github.com/alibug/go-identity-utils/mongoconn"
	"github.com/alibug/go-identity-utils/redisconn"
	"github.com/gin-gonic/gin"
)

func main() {

	duration := config.ReadCustomIntConfig("mongo.duration", false)
	timeDuration := time.Duration(duration) * time.Second

	mongourl := config.ReadMongoConfig("mongo")

	// 3、初始化 MongoDB 数据读取器
	conn, err := mongoconn.NewConn(mongourl, timeDuration)
	if err != nil {
		log.Fatalf("创建数据库连接失败: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = conn.Ping()
	if err != nil {
		log.Fatal("连接MongoDB数据库失败")
	}

	// redisURL := config.ReadRedisConfig("redis")
	// redisConn, err := redisconn.NewConn(redisURL)
	redisConn := redisconn.NewConnFromConfig("redis")
	if err != nil {
		log.Fatalf("创建Redis数据库连接失败:%v", err)
	}

	_, err = redisConn.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatal("连接Redis数据库失败")
	}

	// 4、配置指定的 Collection
	usersColl := conn.GetColl("users")
	userRepo := _userRepo.NewMongoUserRepository(usersColl)
	userUsercase := _userUseCase.NewUserUsecase(userRepo, timeDuration)

	// 5、配置 TokenUserCase
	tokenConfig := config.ReadTokenConfig("token", "maxage")
	tokenRepo := _tokenRepo.NewRedisTokenRepository(redisConn, tokenConfig)
	tokenUsercase := _tokenUseCase.NewTokenUsecase(tokenRepo, timeDuration)

	route := gin.Default()

	cookieConfig := config.ReadCookieConfig("cookie", "maxage")
	_userHttpDelivery.NewUserHandler(route, userUsercase, tokenUsercase, cookieConfig)

	port := config.ReadCustomStringConfig("rest.port")

	route.Run(fmt.Sprintf(":%s", port))
}
