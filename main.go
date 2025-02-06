package main

import (
	"fmt"
	"time"

	initAll "wms-service/init"
	"wms-service/routes"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
)

var db *postgres.DbCluster

func main() {
	// Initialize config -> It will read config.yaml file
	err := config.Init(time.Second * 10)
	if err != nil {
		log.Panicf("Error while initialising config, err: %v", err)
		panic(err)
	}

	ctx, err := config.TODOContext()
	if err != nil {
		fmt.Println("Error creating context.")
	}

	// -----------------------------------------------------------------------------------------------------------------------------------------------------------------

	// Inititalize Database and Logger
	db = initAll.Initialize(ctx)

	// -----------------------------------------------------------------------------------------------------------------------------------------------------------------

	// Initialize Server
	server := http.InitializeServer(config.GetString(ctx, "server.port"), 10*time.Second, 10*time.Second, 70*time.Second)
	log.Infof("Starting server on port" + config.GetString(ctx, "server.port"))

	// Initialize Routes
	errr := routes.Initialize(ctx, server, db)
	if errr != nil {
		log.Errorf(err.Error())
		panic(err)
	}

	// Start Server
	err = server.StartServer("wms-service")
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
}
