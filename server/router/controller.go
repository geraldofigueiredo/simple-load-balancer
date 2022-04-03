package router

import (
	"net/http"
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
	peer := c.serverPool.GetNextPeer()
	if peer == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	peer.ReverseProxy.ServeHTTP(w, r)
}
