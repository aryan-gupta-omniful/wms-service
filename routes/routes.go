package routes

import (
	"context"
	"wms-service/controllers"

	// "github.com/newrelic/go-agent/v3/integrations/nrgin"

	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	// "github.com/omniful/go_commons/newrelic"
)

func Initialize(ctx context.Context, s *http.Server, db *postgres.DbCluster) error {
	// Health Check Route
	s.GET("/health", controllers.GetHealth(db))

	// API v1 Routes Group
	v1 := s.Engine.Group("/api/v1")
	{
		// Hubs Routes
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
