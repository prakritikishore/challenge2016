package main

import (
	"distributor-service/routes"
	"distributor-service/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.properties")
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}
	router := gin.Default()
	router.RedirectTrailingSlash = false
	routes.SetupRoutes(router)
	err = utils.LoadCityData(os.Getenv("app.cityDataPath"))
	if err != nil {
		panic(err)
	}
	router.Run(":" + os.Getenv("app.port"))

}
