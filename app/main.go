package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alibug/go-identity/helper/mongoconn"
	_userHttpDelivery "github.com/alibug/go-identity/user/delivery/restgin"
	_userRepo "github.com/alibug/go-identity/user/repository/mongodb"
	_userUseCase "github.com/alibug/go-identity/user/usecase"
	"github.com/gin-gonic/gin"
)

const mongoURLTemplate = "mongodb://%s:%s@%s/%s"

func main() {
	timeoutDuration := 100 * time.Second

	mongourl := fmt.Sprintf(mongoURLTemplate, mongoUser, mongoPass, mongoHost, mongoDbName)
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

	// 4、配置指定的 Collection
	usersColl := conn.GetColl("users")

	userRepo := _userRepo.NewMongoUserRepository(usersColl)

	userUsercase := _userUseCase.NewUserUsecase(userRepo, timeoutDuration)

	route := gin.Default()
	_userHttpDelivery.NewUserHandler(route, userUsercase)

	route.Run()
}
