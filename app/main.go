package main

import (
	"context"
	"log"
	"time"

	"github.com/alibug/go-identity/helper/configreader"
	"github.com/alibug/go-identity/helper/mongoconn"
	"github.com/alibug/go-identity/helper/redisconn"
	_tokenRepo "github.com/alibug/go-identity/token/repository/redisdb"
	_tokenUseCase "github.com/alibug/go-identity/token/usecase"
	_userHttpDelivery "github.com/alibug/go-identity/user/delivery/restgin"
	_userRepo "github.com/alibug/go-identity/user/repository/mongodb"
	_userUseCase "github.com/alibug/go-identity/user/usecase"
	"github.com/gin-gonic/gin"
)

const mongoURLTemplate = "mongodb://%s:%s@%s/%s"

func main() {

	/*
		mongoUser := "admin017"
		mongoPass := "pass#_27"
		mongoHost := "127.0.0.1:27017"
		mongoDbName := "test_users"
	*/
	timeoutDuration := 100 * time.Second

	mongourl := configreader.ReadMongoConfig()

	// mongourl := fmt.Sprintf(mongoURLTemplate, mongoUser, mongoPass, mongoHost, mongoDbName)
	// 3、初始化 MongoDB 数据读取器
	conn, err := mongoconn.NewConn(mongourl, timeoutDuration)
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

	redisURL := configreader.ReadRedisConfig()
	redisConn, err := redisconn.NewConn(redisURL)
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
	userUsercase := _userUseCase.NewUserUsecase(userRepo, timeoutDuration)

	// 5、配置 TokenUserCase
	tokenConfig := configreader.ReadTokenConfig()
	tokenRepo := _tokenRepo.NewRedisTokenRepository(redisConn, tokenConfig)
	tokenUsercase := _tokenUseCase.NewTokenUsecase(tokenRepo, timeoutDuration)

	route := gin.Default()

	cookieConfig := configreader.ReadCookieConfig()
	_userHttpDelivery.NewUserHandler(route, userUsercase, tokenUsercase, cookieConfig)

	route.Run()
}
