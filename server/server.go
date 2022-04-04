package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	context2 "simple-load-balancer/context"
	"simple-load-balancer/entity"
	"simple-load-balancer/server/router"
	"simple-load-balancer/service"
	"time"
)

const (
	port         = ":8080"
	Attempts int = iota
	Retry
)

var backends = []string{"localhost:8081", "localhost:8082", "localhost:8083"}

func InitServer() {
	serverPool := entity.NewEmptyServerPool()
	serverPoolService := service.NewServerPool(serverPool)
	lbController := router.NewLoadBalancerController(serverPoolService)

	for idx, backendUrlStr := range backends {
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

	_ = http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(lbController.LB),
	}
}
