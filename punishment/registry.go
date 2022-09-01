package punishment

import (
	"fmt"
	"sync"
)

// Container is any data type that can hold user specific data
type Container interface {
	// Data returns a marshable representation of the data without a mutex.
	Data() Dataer
}

// Dataer is returned by a Container for saving and loading. We use a container as it's difficult to directly pass a
// Container with a built in mutex.
type Dataer interface {
	Container() Container
}

// Registry is the base type used to interact with punishments.
type Registry struct {
	provider Provider

	// punishments is a registry of all active punishments on a server. It is first index by puishment type which
	// is represents as map[string]*Container. The identifiers for that map are the specific punishment identifier for
	// the user.
	punishments map[string]map[any]Container

	lock sync.RWMutex
}

// New returns a new punishment handler.
func New(provider Provider) Registry {
	return Registry{
		provider: provider,
	}
}

func (r *Registry) AddAlias(username, ip, xuid string, data ...any) bool {
	ipco, err := r.Load("ip", ip)
	if err != nil {
		return false
	}

	ipc, ok := (*ipco).(*Ip)
	if ok {
		ipc.AddAlias(Alias{})
	}
}

func (r *Registry) Load(ptype string, identifier string) (*Container, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.punishments[ptype]; ok {
		r.punishments[ptype] = map[any]Container{}
	}

	if c, ok := r.punishments[ptype][identifier]; ok {
		return c, nil
	}
	data, err := r.provider.Load(ptype, identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to load container: %w", err)
	}
	container := data.Container()

	r.punishments[ptype][identifier] = container
	return &r.punishments[ptype][identifier], nil
}

// Save attempts to save all the data within the punishment Registry.
func (r *Registry) Save() error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	var err error
	err = nil
	for ptype, punishments := range r.punishments {
		for id, punishment := range punishments {
			p := *punishment
			data := p.Data()
			er := r.provider.Save(ptype, id, data)
			if err == nil && er != nil {
				err = fmt.Errorf("error saving punishment type: %v identifier %v: %w", ptype, id, er)
			}
		}
	}
	return err
}

// Close ...
func (r *Registry) Close() error {
	return r.Save()
}
