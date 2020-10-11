package main

import (
	"dogo/pkg/api"
	"dogo/pkg/config"
	"dogo/pkg/model"
	"dogo/pkg/repository"
	"dogo/pkg/utils"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func main() {

	config.Dogo = config.SetupConfig()

	var err error
	config.DB, err = gorm.Open(mysql.Open(config.Dogo.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("连接数据库异常", err)
	}

	if err := config.DB.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}

	var users = make([]model.User, 0)
	if err := repository.FindAllUser(&users); err != nil {
		panic(err)
	}

	if len(users) == 0 {

		var pass []byte
		if pass, err = utils.Encoder.Encode([]byte("admin")); err != nil {
			return
		}

		user := model.User{
			ID:       utils.UUID(),
			Username: "admin",
			Password: string(pass),
			Nickname: "超级管理员",
			Type:     "manager",
			Created:  utils.NowJsonTime(),
		}
		if err := repository.CreateNewUser(&user); err != nil {
			panic(err)
		}
	}

	if err := config.DB.AutoMigrate(&model.Asset{}); err != nil {
		panic(err)
	}
	if err := config.DB.AutoMigrate(&model.Session{}); err != nil {
		panic(err)
	}
	if err := config.DB.AutoMigrate(&model.Command{}); err != nil {
		panic(err)
	}
	if err := config.DB.AutoMigrate(&model.Credential{}); err != nil {
		panic(err)
	}

	config.Cache = cache.New(5*time.Minute, 10*time.Minute)
	r := api.SetupRoutes()

	if err := r.Run(config.Dogo.Addr...); err != nil {
		panic(err)
	}
}
