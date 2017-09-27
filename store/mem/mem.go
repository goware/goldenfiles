package mem

import (
	"errors"
	"sync"

	"github.com/goware/mockingbird/store"
)

type Mem struct {
	data map[string][]byte
	mu   sync.Mutex
}

func (m *Mem) Get(key string) (buf []byte, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data != nil {
		if v, ok := m.data[key]; ok {
			return v, nil
		}
	}
	return nil, errors.New("no such value")
}

func (m *Mem) Set(key string, buf []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = map[string][]byte{}
	}
	m.data[key] = buf
	return nil
}

func (m *Mem) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data != nil {
		delete(m.data, key)
	}
	return nil
}

var _ = store.Store(&Mem{})
