package services

import (
	"distributor-service/models"
	"distributor-service/stores"
	"distributor-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func AddDistributor(ctx *gin.Context) {
	var distributor models.Distributor
	if err := ctx.ShouldBindJSON(&distributor); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if _, exists := stores.GetDistributors()[distributor.Name]; exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Distributor already Present"})
	} else {
		distributor.Permissions = make(map[string]bool)
		distributor.SubDistributors = make(map[string]bool)
		distributor.ParentDistributors = make(map[string]bool)
		stores.SaveDistributor(&distributor)
		distributor, _ := stores.GetDistributorByName(distributor.Name)
		ctx.JSON(http.StatusCreated, distributor)
	}
}

func GetDistributor(ctx *gin.Context) {
	distributorName := ctx.Param("distributor-name")
	if distributor, exists := stores.GetDistributorByName(distributorName); exists {
		ctx.JSON(http.StatusOK, distributor)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Distributor not found"})
	}
}

func DeleteDistributor(ctx *gin.Context) {
	distributorName := ctx.Param("distributor-name")
	if _, exists := stores.GetDistributorByName(distributorName); exists {
		stores.DeleteDistributor(distributorName)
		ctx.JSON(http.StatusNoContent, true)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Distributor not found"})
	}
}

func AddPermission(ctx *gin.Context) {
	var req models.AddPermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	distributorName := ctx.Param("distributor-name")
	distributor, exists := stores.GetDistributorByName(distributorName)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Distributor not found"})
		return
	}

	if _, valid := stores.RegionCodeToNameLookup[req.RegionCode]; valid {
		distributor.Permissions[req.RegionCode] = req.IsIncluded
		stores.SaveDistributor(distributor)
		ctx.JSON(http.StatusOK, distributor)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region code provided"})
	}
}

func DeletePermission(ctx *gin.Context) {
	var req models.DeletePermissionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	distributorName := ctx.Param("distributor-name")
	distributor, exists := stores.GetDistributorByName(distributorName)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Distributor not found"})
		return
	}

	delete(distributor.Permissions, req.RegionCode)
	stores.SaveDistributor(distributor)
	ctx.JSON(http.StatusOK, distributor)
}

func CheckPermission(ctx *gin.Context) {
	var req models.CheckPermissionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	distributor, exists := stores.GetDistributorByName(req.Name)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Distributor not found"})
		return
	}

	if utils.HasPermission(distributor, req.RegionCode) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Allowed to distribute"})
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Not allowed to distribute"})
	}
}

func AuthorizeSubDistributor(ctx *gin.Context) {
	var req models.AuthorizeDistributorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from, existsFrom := stores.GetDistributorByName(req.FromDistributor)
	to, existsTo := stores.GetDistributorByName(req.ToDistributor)
	if !existsFrom || !existsTo {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Distributor not found"})
		return
	}

	if !utils.IsSubDistributor(to, from) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Authorization failed: Not a sub-distributor"})
		return
	}

	if _, valid := stores.RegionCodeToNameLookup[req.RegionCode]; !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region code provided"})
	}

	if utils.HasPermission(from, req.RegionCode) {
		to.Permissions[req.RegionCode] = true
		ctx.JSON(http.StatusOK, gin.H{"message": "Distributor authorized successfully"})
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Authorization failed: insufficient permissions"})
	}
}

func AddSubDistributor(ctx *gin.Context) {
	var req models.SubDistributorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	distributor, exists := stores.GetDistributorByName(req.DistributorName)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Parent distributor not found"})
		return
	}

	if _, exists := stores.GetDistributorByName(req.SubDistributorName); exists {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Sub-distributor already exists"})
		return
	}

	subDistributor := &models.Distributor{
		Name:        req.SubDistributorName,
		Permissions: make(map[string]bool),
	}
	subDistributor.ParentDistributors = make(map[string]bool)
	subDistributor.ParentDistributors[distributor.Name] = true
	stores.SaveDistributor(subDistributor)
	distributor.SubDistributors[req.SubDistributorName] = true
	stores.SaveDistributor(distributor)

	ctx.JSON(http.StatusCreated, subDistributor)
}

func DeleteSubDistributor(ctx *gin.Context) {
	distributorName := ctx.Param("distributor-name")
	subDistributorName := ctx.Param("sub-distributor-name")

	_, subExists := stores.GetDistributorByName(subDistributorName)
	if !subExists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Sub-distributor distributor not found"})
		return
	}

	_, exists := stores.GetDistributorByName(distributorName)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Parent distributor not found"})
		return
	}

	stores.DeleteSubDistributor(distributorName, subDistributorName)
	stores.DeleteDistributor(subDistributorName)

	ctx.JSON(http.StatusOK, gin.H{"message": "Sub-distributor deleted successfully"})
}
