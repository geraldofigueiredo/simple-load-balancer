package router

import (
	"log"
	"net/http"
	context2 "simple-load-balancer/context"
	"simple-load-balancer/contract"
)

type LoadBalancerController struct {
	serverPool contract.ServerPoolService
}

func NewLoadBalancerController(serverPool contract.ServerPoolService) *LoadBalancerController {
	return &LoadBalancerController{
		serverPool: serverPool,
	}
}

func (c *LoadBalancerController) LB(w http.ResponseWriter, r *http.Request) {
	attempts := context2.GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := c.serverPool.GetNextPeer()
	if peer == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	peer.ReverseProxy.ServeHTTP(w, r)
}
