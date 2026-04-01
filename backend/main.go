package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"restaurant-management-system/config"
	"restaurant-management-system/db"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := db.InitPostgres(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.ClosePostgres()

	router := gin.Default()
	v1 := router.Group("/v1")
	defineRoutes(v1)

	addr := fmt.Sprintf(":%d", config.RuntimeConfig.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
