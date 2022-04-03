package server

import (
	"net/http"
	"simple-load-balancer/entity"
	"simple-load-balancer/server/router"
	"simple-load-balancer/service"
)

func InitServer() {
	serverPool := entity.NewEmptyServerPool()
	serverPool.AddBackend(entity.NewBackend("backend-1", "http://localhost:8081"))
	serverPool.AddBackend(entity.NewBackend("backend-2", "http://localhost:8082"))
	serverPool.AddBackend(entity.NewBackend("backend-3", "http://localhost:8083"))

	serverPoolService := service.NewServerPool(serverPool)

	lbController := router.NewLoadBalancerController(serverPoolService)

	_ = http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(lbController.LB),
	}
}
