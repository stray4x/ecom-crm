package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stray4x/ecom-crm/internal/config"
	db "github.com/stray4x/ecom-crm/internal/database"
)

func main() {
	db.NewDB(config.Env)

	router := gin.Default()

	log.Fatal(router.Run(":" + config.Env.AppPort))
}
