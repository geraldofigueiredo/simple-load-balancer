package service

import (
	"simple-load-balancer/contract"
	"simple-load-balancer/entity"
	"sync/atomic"
)

type serverPoolService struct {
	pool *entity.ServerPool
}

func NewServerPool(serverPool *entity.ServerPool) contract.ServerPoolService {
	return &serverPoolService{
		pool: serverPool,
	}
}

func (s *serverPoolService) NextIndex() int {
	return int(atomic.AddUint64(&s.pool.Current, uint64(1)) % uint64(len(s.pool.Backends)))
}

func (s *serverPoolService) GetNextPeer() *entity.Backend {
	next := s.NextIndex()
	l := len(s.pool.Backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.pool.Backends)
		if s.pool.Backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.pool.Current, uint64(idx))
			}
			return s.pool.Backends[idx]
		}
	}
	return nil
}
