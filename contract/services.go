package contract

import "simple-load-balancer/entity"

type ServerPoolService interface {
	NextIndex() int
	GetNextPeer() *entity.Backend
}
