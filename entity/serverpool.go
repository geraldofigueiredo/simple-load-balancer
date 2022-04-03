package entity

type ServerPool struct {
	Backends []*Backend
	Current  uint64
}

func NewEmptyServerPool() *ServerPool {
	return &ServerPool{
		Backends: make([]*Backend, 0),
		Current:  0,
	}
}

func (s *ServerPool) AddBackend(backend *Backend) {
	s.Backends = append(s.Backends, backend)
}
