package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stray4x/ecom-crm/internal/config"
	"github.com/stray4x/ecom-crm/internal/database"
	"github.com/stray4x/ecom-crm/internal/handlers"
	"github.com/stray4x/ecom-crm/internal/repository"
	"github.com/stray4x/ecom-crm/internal/routes"
	"github.com/stray4x/ecom-crm/internal/service"
)

func main() {
	db := database.NewDB(config.Env)

	customerRepo := repository.NewCustomerRepository(db)

	authService := service.NewAuthService(customerRepo)

	handlers := &routes.Handlers{
		Auth: handlers.NewAuthHandler(authService),
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	routes.Setup(router, handlers)

	log.Fatal(router.Run(":" + config.Env.AppPort))
}
