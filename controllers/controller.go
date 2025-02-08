package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"wms-service/models"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/db/sql/postgres"
)

type ValidateOrderRequest struct {
	SKUID string `json:"sku_id"`
	HubID string `json:"hub_id"`
}

func GetHealth(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("server is workign absolutely fine - OK")
	}
}

func GetAllHubs(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var hubs []models.Hub
		result := db.GetMasterDB(ctx).Find(&hubs)
		if result.Error != nil {
			ctx.JSON(500, gin.H{"error": "Failed to fetch hubs", "details": result.Error.Error()})
			return
		}
		if len(hubs) == 0 {
			ctx.JSON(404, gin.H{"message": "No hubs found"})
			return
		}
		ctx.JSON(200, gin.H{"hubs": hubs})
	}
}

func GetHubByID(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(400, gin.H{"error": "ID parameter is required"})
			return
		}

		var hub models.Hub
		result := db.GetMasterDB(ctx).First(&hub, id)
		if result.Error != nil {
			ctx.JSON(404, gin.H{"error": "Hub not found"})
			return
		}
		ctx.JSON(200, gin.H{"hub": hub})
	}
}

func CreateHub(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var newHub models.Hub
		if err := ctx.BindJSON(&newHub); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request payload", "details": err.Error()})
			return
		}

		// Basic validation
		if newHub.ManagerEmail == "" || newHub.ManagerName == "" {
			ctx.JSON(400, gin.H{"error": "Manager email and name are required"})
			return
		}

		if err := db.GetMasterDB(ctx).Create(&newHub).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create hub", "details": err.Error()})
			return
		}

		ctx.JSON(201, gin.H{"message": "Hub created successfully", "hub": newHub})
	}
}

/* {
	"id": 11,
	"tenant_id": 6,
	"manager_name": "aryan",
	"manager_contact": "9876019282",
	"manager_email": "aryan@email.com"
} */

func GetAllSKUs(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var skus []models.SKU
		result := db.GetMasterDB(ctx).Find(&skus)
		if result.Error != nil {
			ctx.JSON(500, gin.H{"error": "Failed to fetch SKUs", "details": result.Error.Error()})
			return
		}
		if len(skus) == 0 {
			ctx.JSON(404, gin.H{"message": "No SKUs found"})
			return
		}
		ctx.JSON(200, gin.H{"skus": skus})
	}
}

func GetSKUByID(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(400, gin.H{"error": "ID parameter is required"})
			return
		}

		var sku models.SKU
		result := db.GetMasterDB(ctx).First(&sku, id)
		if result.Error != nil {
			ctx.JSON(404, gin.H{"error": "SKU not found"})
			return
		}
		ctx.JSON(200, gin.H{"sku": sku})
	}
}

func CreateSKU(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var newSKU models.SKU
		if err := ctx.BindJSON(&newSKU); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request payload", "details": err.Error()})
			return
		}

		// Add any SKU-specific validation here
		// Example (adjust according to your SKU model requirements):
		if newSKU.ProductID == 0 {
			ctx.JSON(400, gin.H{"error": "Product ID is required"})
			return
		}

		if err := db.GetMasterDB(ctx).Create(&newSKU).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create SKU", "details": err.Error()})
			return
		}

		ctx.JSON(201, gin.H{"message": "SKU created successfully", "sku": newSKU})
	}
}

type ValidationResponse struct {
	IsValid bool
	Error   error
}

func ValidateHubAndSKU(db *postgres.DbCluster) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("ValidateHubAndSKU function called")
		// ch := make(chan string, 1)

		// http://localhost:8082/api/v1/hubs/:id
		// http://localhost:8082/api/v1/skus/:id

		var request ValidateOrderRequest

		// Bind JSON request body to struct
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request payload",
				"details": err.Error(),
			})
			return
		}

		log.Println("parsed request from ValidateHubAndSKU function is: ", request)

		// Now you can access the values using:
		// request.SKUID
		// request.HubID

		// Your validation logic here...
		var wg sync.WaitGroup
		respChan := make(chan ValidationResponse, 2)

		// Start hub validation goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			var hub models.Hub
			result := db.GetMasterDB(ctx).First(&hub, request.HubID)
			if result.Error != nil {
				respChan <- ValidationResponse{
					IsValid: false,
					Error:   fmt.Errorf("hub validation failed: %v", result.Error),
				}
				return
			}
			respChan <- ValidationResponse{IsValid: true}
		}()

		// Start SKU validation goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			var sku models.SKU
			result := db.GetMasterDB(ctx).First(&sku, request.SKUID)
			if result.Error != nil {
				respChan <- ValidationResponse{
					IsValid: false,
					Error:   fmt.Errorf("SKU validation failed: %v", result.Error),
				}
				return
			}
			respChan <- ValidationResponse{IsValid: true}
		}()

		// Wait for both goroutines to complete
		go func() {
			wg.Wait()
			close(respChan)
		}()

		// Process results
		var finalResponse ValidationResponse
		finalResponse.IsValid = true

		for resp := range respChan {
			if !resp.IsValid {
				finalResponse.IsValid = false
				finalResponse.Error = resp.Error
				break
			}
		}

		// Send response based on validation result
		if finalResponse.IsValid {
			ctx.JSON(http.StatusOK, gin.H{"message": "Validation successful"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": finalResponse.Error.Error()})
		}
	}
}
