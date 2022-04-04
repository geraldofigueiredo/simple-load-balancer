package entity

type ServerPool struct {
	Backends []*Backend
	Current  uint64
	count    int
}

// NewEmptyServerPool creates a new ServerPool
// count == -1 means empty pool
func NewEmptyServerPool() *ServerPool {
	return &ServerPool{
		Backends: make([]*Backend, 0),
		Current:  0,
		count:    -1,
	}
}

func (s *ServerPool) AddBackend(backend *Backend) (index int) {
	s.Backends = append(s.Backends, backend)
	s.count++
	return s.count
}

func (s *ServerPool) GetPoolSize() int {
	return s.count
}

func (s *ServerPool) SetBackendAlive(index int, alive bool) {
	s.Backends[index].SetAlive(alive)
}

func (s *ServerPool) GetBackend(index int) *Backend {
	return s.Backends[index]
}
