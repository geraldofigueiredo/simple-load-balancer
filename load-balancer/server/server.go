package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	context2 "simple-load-balancer/context"
	"simple-load-balancer/entity"
	"simple-load-balancer/server/router"
	"simple-load-balancer/service"
	"strings"
	"time"
)

const (
	port = 8080
)

func InitServer() {
	backends := os.Getenv("BACKENDS")
	if backends == "" {
		backends = "http://localhost:8081,http://localhost:8082,http://localhost:8083"
	}
	serverPool := entity.NewEmptyServerPool()
	serverPoolService := service.NewServerPool(serverPool)
	lbController := router.NewLoadBalancerController(serverPoolService)

	for idx, backendUrlStr := range strings.Split(backends, ",") {
		serverUrl, err := url.Parse(backendUrlStr)
		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(serverUrl)

		backendIndex := serverPool.AddBackend(entity.NewBackend(fmt.Sprintf("backend-%d", idx), backendUrlStr, proxy))

		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := context2.GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), context2.Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this backend as down
			serverPool.SetBackendAlive(backendIndex, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := context2.GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), context2.Attempts, attempts+1)
			lbController.LB(writer, request.WithContext(ctx))
		}
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lbController.LB),
	}

	go serverPoolService.HealthCheck()

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
