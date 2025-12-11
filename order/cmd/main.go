package main

import (
	"log"

	"github.com/pauloabaia/microservices/order/config"
	"github.com/pauloabaia/microservices/order/internal/adapters/db"
	"github.com/pauloabaia/microservices/order/internal/adapters/grpc"
	"github.com/pauloabaia/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter)
	grpcSrv := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcSrv.Run()
}
