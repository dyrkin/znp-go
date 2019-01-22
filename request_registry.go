package znp

import (
	"sync"
	"time"

	unp "github.com/dyrkin/unp-go"
)

type registryKey struct {
	subsystem unp.Subsystem
	command   byte
}
type registryValue struct {
	syncRsp  chan *unp.Frame
	syncErr  chan error
	deadline *deadline
}

type deadline struct {
	timer     *time.Timer
	cancelled chan bool
}

func (d *deadline) Cancel() {
	d.cancelled <- true
}

type RequestRegistry struct {
	registry map[registryKey]*registryValue
	m        sync.RWMutex
}

func NewRequestRegistry() *RequestRegistry {
	return &RequestRegistry{registry: map[registryKey]*registryValue{}}
}

func (r *RequestRegistry) Register(key *registryKey, value *registryValue) {
	r.m.Lock()
	defer r.m.Unlock()
	r.registry[*key] = value
}

func (r *RequestRegistry) Unregister(key *registryKey) {
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.registry, *key)
}

func (r *RequestRegistry) Get(key *registryKey) (*registryValue, bool) {
	r.m.RLock()
	r.m.RUnlock()
	value, ok := r.registry[*key]
	return value, ok
}
