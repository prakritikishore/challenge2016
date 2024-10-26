package tests

import (
	"bytes"
	"distributor-service/models"
	"distributor-service/routes"
	"distributor-service/stores"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.SetupRoutes(router)
	return router
}

func TestAddDistributor(t *testing.T) {
	router := setupRouter()

	_ = models.Distributor{Name: "Distributor1"}

	body := bytes.NewBufferString(`{"name":"Distributor1"}`)
	req, _ := http.NewRequest("POST", "/distributor-service/distributor/add", body)
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Contains(t, writer.Body.String(), "Distributor1")
}

func TestGetDistributor(t *testing.T) {
	router := setupRouter()

	distributor := models.Distributor{Name: "Distributor1"}
	stores.SaveDistributor(&distributor)

	req, _ := http.NewRequest("GET", "/distributor-service/distributor/Distributor1", nil)
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "Distributor1")
}

func TestDeleteDistributor(t *testing.T) {
	router := setupRouter()
	distributor := models.Distributor{Name: "Distributor1"}
	stores.SaveDistributor(&distributor)

	req, _ := http.NewRequest("DELETE", "/distributor-service/distributor/Distributor1", nil)
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusNoContent, writer.Code)

	_, exists := stores.GetDistributorByName("Distributor1")
	assert.False(t, exists)
}

func TestAddSubDistributor(t *testing.T) {
	router := setupRouter()

	parentDistributor := models.Distributor{Name: "Distributor1", Permissions: make(map[string]bool)}
	parentDistributor.SubDistributors = make(map[string]bool)
	stores.SaveDistributor(&parentDistributor)
	body := bytes.NewBufferString(`{"distributor_name":"Distributor1", "sub_distributor_name":"Distributor2"}`)
	req, _ := http.NewRequest("POST", "/distributor-service/sub-distributor/add", body)
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Contains(t, writer.Body.String(), "Distributor2")
}

func TestDeleteSubDistributor(t *testing.T) {
	router := setupRouter()

	parentDistributor := models.Distributor{Name: "Distributor1", Permissions: make(map[string]bool)}
	parentDistributor.SubDistributors = map[string]bool{"Distributor2": true}
	stores.SaveDistributor(&parentDistributor)

	subDistributor := models.Distributor{Name: "Distributor2", Permissions: make(map[string]bool)}
	subDistributor.ParentDistributors = map[string]bool{"Distributor1": true}
	stores.SaveDistributor(&subDistributor)
	stores.SaveDistributor(&parentDistributor)
	req, _ := http.NewRequest("DELETE", "/distributor-service/sub-distributor/Distributor1/Distributor2", nil)
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusOK, writer.Code)
	fmt.Println(writer.Body.String())
	assert.Contains(t, writer.Body.String(), "Sub-distributor deleted successfully")

	_, exists := stores.GetDistributorByName("Distributor2")
	assert.False(t, exists)
}

func TestAddPermission(t *testing.T) {
	router := setupRouter()

	distributor := models.Distributor{Name: "Distributor1", Permissions: map[string]bool{}}
	stores.SaveDistributor(&distributor)
	stores.RegionCodeToNameLookup = map[string]string{"IN": "IN"}
	reqBody := `{"region_code": "IN", "is_included": true}`
	req, _ := http.NewRequest("POST", "/distributor-service/permissions/Distributor1", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "Distributor1")
}

func TestDeletePermission(t *testing.T) {
	router := setupRouter()

	distributor := models.Distributor{Name: "Distributor1", Permissions: map[string]bool{"IN": true}}
	stores.SaveDistributor(&distributor)

	reqBody := `{"region_code": "IN"}`
	req, _ := http.NewRequest("DELETE", "/distributor-service/permissions/Distributor1", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "Distributor1")
}

func TestCheckPermission(t *testing.T) {
	router := setupRouter()

	distributor := models.Distributor{Name: "Distributor1", Permissions: map[string]bool{"IN": true}}
	stores.SaveDistributor(&distributor)

	reqBody := `{"name": "Distributor1", "region_code": "IN"}`
	req, _ := http.NewRequest("POST", "/distributor-service/check-permissions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "Allowed to distribute")
}

func TestAuthorizeSubDistributor(t *testing.T) {
	router := setupRouter()

	fromDistributor := models.Distributor{Name: "Distributor1", Permissions: map[string]bool{"IN": true}}
	fromDistributor.SubDistributors = map[string]bool{"Distributor2": true}
	toDistributor := models.Distributor{Name: "Distributor2", Permissions: map[string]bool{}}
	stores.SaveDistributor(&fromDistributor)
	stores.SaveDistributor(&toDistributor)

	reqBody := `{"from_distributor": "Distributor1", "to_distributor": "Distributor2", "region_code": "IN"}`
	req, _ := http.NewRequest("POST", "/distributor-service/authorize/sub-distributor", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, req)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "Distributor authorized successfully")
}
