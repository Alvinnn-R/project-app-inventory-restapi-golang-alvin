package main

import (
	"fmt"
	"log"
	"net/http"
	"project-app-inventory/database"
	"project-app-inventory/handler"
	"project-app-inventory/repository"
	"project-app-inventory/router"
	"project-app-inventory/service"
	"project-app-inventory/utils"
)

func main() {
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatal("error file configration")
	}
	fmt.Println(config)

	// Initialize logger FIRST
	logger, err := utils.InitLogger(config.PathLogging, config.Debug)
	if err != nil {
		log.Fatal("error initializing logger: ", err)
	}

	db, err := database.InitDB(logger, config.DB)
	if err != nil {
		panic(err)
	}

	// logger, err := utils.InitLogger(config.PathLogging, config.Debug)

	repo := repository.NewRepository(db, logger)
	service := service.NewService(repo)
	handler := handler.NewHandler(service, config)

	r := router.NewRouter(handler, service, logger)

	fmt.Println("server running on port " + config.Port)
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Fatal("error server: ", err)
	}
}
