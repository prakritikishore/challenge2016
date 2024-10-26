package routes

import (
	"distributor-service/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/health", services.HealthCheck)
	distributorGroup := router.Group("/distributor-service")
	{
		distributorGroup.POST("/distributor/add", services.AddDistributor)
		distributorGroup.GET("/distributor/:distributor-name", services.GetDistributor)
		distributorGroup.DELETE("/distributor/:distributor-name", services.DeleteDistributor)

		distributorGroup.POST("/sub-distributor/add", services.AddSubDistributor)
		distributorGroup.DELETE("/sub-distributor/:distributor-name/:sub-distributor-name", services.DeleteSubDistributor)

		distributorGroup.POST("/permissions/:distributor-name", services.AddPermission)
		distributorGroup.DELETE("/permissions/:distributor-name", services.DeletePermission)
		distributorGroup.POST("/check-permissions", services.CheckPermission)

		distributorGroup.POST("/authorize/sub-distributor", services.AuthorizeSubDistributor)
	}
}
