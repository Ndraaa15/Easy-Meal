package main

import (
	"bcc-project-v/sdk/config"
	"bcc-project-v/sdk/database"
	"bcc-project-v/src/entities"
	"bcc-project-v/src/handlers"
	"bcc-project-v/src/repository"
	"fmt"
)

func main() {
	conf := config.Init()
	if err := conf.LoadEnv(".env"); err != nil {
		panic(err)
	}

	sqlConfig := database.Config{
		Username: conf.GetEnv("DB_USERNAME"),
		Password: conf.GetEnv("DB_PASSWORD"),
		Host:     conf.GetEnv("DB_HOST"),
		Port:     conf.GetEnv("DB_PORT"),
		Database: conf.GetEnv("DB_DATABASE"),
	}

	db, err := database.Init(sqlConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully!")

	db.Debug().AutoMigrate(&entities.User{}, &entities.Seller{})

	handler := handlers.Init(conf, repository.NewRepository(db))
	handler.Run()

}
