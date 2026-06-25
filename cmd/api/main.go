package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stray4x/ecom-crm/internal/config"
	"github.com/stray4x/ecom-crm/internal/database"
	"github.com/stray4x/ecom-crm/internal/handlers"
	"github.com/stray4x/ecom-crm/internal/redis"
	"github.com/stray4x/ecom-crm/internal/repository"
	"github.com/stray4x/ecom-crm/internal/routes"
	"github.com/stray4x/ecom-crm/internal/service"
)

func main() {
	cfg := config.InitConfig()
	db := database.NewDB(cfg)
	rdb, err := redis.NewClient(cfg.RedisHost, cfg.RedisPort)

	if err != nil {
		log.Fatal(err)
	}

	tokenStore := redis.NewTokenStore(rdb)

	customerRepo := repository.NewCustomerRepository(db)

	authService := service.NewAuthService(customerRepo, cfg, tokenStore)

	handlers := &routes.Handlers{
		Auth: handlers.NewAuthHandler(authService, cfg),
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	routes.Setup(router, handlers)

	log.Fatal(router.Run(":" + cfg.AppPort))
}
