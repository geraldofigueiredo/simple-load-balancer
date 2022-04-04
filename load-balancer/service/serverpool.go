package service

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"simple-load-balancer/contract"
	"simple-load-balancer/entity"
	"sync/atomic"
	"time"
)

const (
	SiteUnreachableMessage = "Site unreachable"
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

func (s *serverPoolService) HealthCheck() {
	for _, b := range s.pool.Backends {
		alive := isBackendAlive(b.URL.String())
		status := getStatusBasedOnAlive(alive)
		b.SetAlive(alive)
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

func getStatusBasedOnAlive(alive bool) string {
	if alive {
		return "up"
	}
	return "down"
}

// isBackendAlive checks if the backend is alive by pinging it
func isBackendAlive(urlStr string) bool {
	u, err := url.Parse(urlStr + "/health")
	if err != nil {
		log.Printf("Error parsing url: %s\n", err)
		return false
	}

	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println(fmt.Sprintf("%s, error: %s", SiteUnreachableMessage, err))
		return false
	}
	defer conn.Close()
	return true
}
