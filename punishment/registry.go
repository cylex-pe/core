package punishment

import (
	"fmt"
	"sync"
)

const IpIdentifier = "ip"
const XuidIdentifier = "xuid"
const DeviceIdentifier = "device"

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

type AliasHandler func(username, ip, device, xuid string, data ...any) bool

// Registry is the base type used to interact with punishments.
type Registry struct {
	provider Provider

	// punishments is a registry of all active punishments on a server. It is first index by puishment type which
	// is represents as map[string]*Container. The identifiers for that map are the specific punishment identifier for
	// the user.
	punishments map[string]map[any]Container
	// aliasHandler is called when a new alias is added with AddAlias.
	aliasHandler AliasHandler
	lock         sync.RWMutex
}

// New returns a new punishment handler.
func New(provider Provider, aliasHandler AliasHandler) Registry {
	return Registry{
		provider:     provider,
		aliasHandler: aliasHandler,
	}
}

// AddAlias will register an alias onto ip and device as well as call the aliasHandler.
func (r *Registry) AddAlias(username, ip, device, xuid string, data ...any) bool {
	ipc, err := r.Ip(ip)
	if err == nil {
		ipc.AddAlias(Alias{
			Username: username,
			Xuid:     xuid,
		})
	}
	dev, err := r.Device(device)
	if err == nil {
		dev.AddAlias(Alias{
			Username: username,
			Xuid:     xuid,
		})
	}
	return r.aliasHandler(username, ip, device, xuid, data)
}

// Xbox attempts to load an xbox object and return it.
func (r *Registry) Xbox(xuid string) (*Xbox, error) {
	xboxc, err := r.Load(XuidIdentifier, xuid)
	if err != nil {
		return nil, err
	}
	xbox, ok := (xboxc).(*Xbox)
	if !ok {
		return nil, fmt.Errorf("container type is not of type Xbox")
	}
	return xbox, nil
}

// Ip attempts to load an ip object and return it.
func (r *Registry) Ip(ip string) (*Ip, error) {
	ipco, err := r.Load(IpIdentifier, ip)
	if err != nil {
		return nil, err
	}
	ipc, ok := (ipco).(*Ip)
	if !ok {
		return nil, fmt.Errorf("container type is not of type Ip")
	}
	return ipc, nil
}

// Device attempts to load a device object and return it.
func (r *Registry) Device(device string) (*Device, error) {
	devc, err := r.Load(DeviceIdentifier, device)
	if err != nil {
		return nil, err
	}
	dev, ok := (devc).(*Device)
	if !ok {
		return nil, fmt.Errorf("container type is not of type Device")
	}
	return dev, nil
}

// Load will attempt to load a Container from the provider and return it. It takes in a punishment type and a user
// identifier.
func (r *Registry) Load(ptype string, identifier any) (Container, error) {
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
	return r.punishments[ptype][identifier], nil
}

// Save attempts to save all the data within the punishment Registry.
func (r *Registry) Save() error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	var err error
	err = nil
	for ptype, punishments := range r.punishments {
		for id, punishment := range punishments {
			p := punishment
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
