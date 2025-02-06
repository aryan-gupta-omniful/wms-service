package controllers

import (
	"fmt"
	"wms-service/models"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/db/sql/postgres"
)

func GetHealth(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("server is workign absolutely fine - OK")
	}
}

func GetAllHubs(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Get all Hubs")
		var check []models.Hub
		db.GetMasterDB(ctx).Find(&check)
		ctx.JSON(200, gin.H{"all hubs": check})
	}
}

func GetHubByID(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		fmt.Println("Get Hub By ID")
		var oneHub models.Hub
		db.GetMasterDB(ctx).Find(&oneHub, id)
		ctx.JSON(200, gin.H{"one hub by id": oneHub})
	}
}

func CreateHub(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var newHub models.Hub
		if err := ctx.BindJSON(&newHub); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := db.GetMasterDB(ctx).Create(&newHub).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create hub"})
			return
		}

		ctx.JSON(200, gin.H{"created_hub": newHub})
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
		fmt.Println("Get all SKUs")
		var allSKUs []models.SKU
		db.GetMasterDB(ctx).Find(&allSKUs)
		ctx.JSON(200, gin.H{"all skus": allSKUs})
	}
}

func GetSKUByID(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		fmt.Println("Get SKU By ID")
		var oneSKU models.SKU
		db.GetMasterDB(ctx).Find(&oneSKU, id)
		ctx.JSON(200, gin.H{"one SKU by id": oneSKU})
	}
}

func CreateSKU(db *postgres.DbCluster) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var newSKU models.SKU
		if err := ctx.BindJSON(&newSKU); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := db.GetMasterDB(ctx).Create(&newSKU).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create sku"})
			return
		}

		ctx.JSON(200, gin.H{"created_sku": newSKU})
	}
}
