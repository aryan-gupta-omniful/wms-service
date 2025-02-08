package routes

import (
	"context"
	"wms-service/controllers"

	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
)

func Initialize(ctx context.Context, s *http.Server, db *postgres.DbCluster) error {
	// Health Check Route
	s.GET("/health", controllers.GetHealth(db))

	// API v1 Routes Group
	v1 := s.Engine.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("/validate_order", controllers.ValidateHubAndSKU(db))
			orders.POST("/validate_inventory", controllers.ValidateAndUpdateInventory(db))
		}

		hubs := v1.Group("/hubs")
		{
			hubs.GET("", controllers.GetAllHubs(db))
			hubs.GET("/:id", controllers.GetHubByID(db))
			hubs.POST("", controllers.CreateHub(db))
		}

		skus := v1.Group("/skus")
		{
			skus.GET("", controllers.GetAllSKUs(db))
			skus.GET("/:id", controllers.GetSKUByID(db))
			skus.POST("", controllers.CreateSKU(db))
		}
	}

	log.Infof("Routes initialized successfully")
	return nil
}
