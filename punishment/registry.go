package punishment

import (
	"fmt"
	"sync"
)

type Registry struct {
	provider Provider

	xboxMu          sync.RWMutex
	xboxPunishments map[string]*Xbox

	ipMu          sync.RWMutex
	ipPunishments map[string]*Ip
}

// New returns a new punishment handler.
func New(provider Provider) Registry {
	return Registry{
		provider: provider,
	}
}

// Xbox attempts to load an xbox xuid punishment holder from the handler or the database.
func (r *Registry) Xbox(xuid string) (*Xbox, error) {
	r.xboxMu.Lock()
	defer r.xboxMu.Unlock()
	if xbox, ok := r.xboxPunishments[xuid]; ok {
		return xbox, nil
	}
	data, err := r.provider.LoadXbox(xuid)
	if err != nil {
		return nil, fmt.Errorf("unable to load holder: %w", err)
	}
	xbox := Xbox{
		username:   data.Username,
		currentBan: data.CurrentBan,
		pastBans:   data.PastBans,
	}
	r.xboxPunishments[xuid] = &xbox
	return &xbox, nil
}

func (r *Registry) Ip(address string) (*Ip, error) {
	r.ipMu.Lock()
	defer r.ipMu.Unlock()
	if ip, ok := r.ipPunishments[address]; ok {
		return ip, nil
	}
	data, err := r.provider.LoadIp(address)
	if err != nil {
		return nil, fmt.Errorf("unable to load ip: %w", err)
	}
	ip := Ip{
		aliases:    data.Aliases,
		currentBan: data.CurrentBan,
		pastBans:   data.PastBans,
	}
	r.ipPunishments[address] = &ip
	return &ip, nil
}

// Save attempts to save all the data within the punishment Registry.
func (r *Registry) Save() error {
	r.xboxMu.RLock()
	r.ipMu.RLock()
	defer r.xboxMu.RUnlock()
	defer r.ipMu.RUnlock()
	var err error
	err = nil
	for xuid, punishment := range r.xboxPunishments {
		data := punishment.Data()
		er := r.provider.SaveXbox(xuid, data)
		if err == nil && er != nil {
			err = fmt.Errorf("error saving punishments: %w", er)
		}
	}
	for ip, punishment := range r.ipPunishments {
		data := punishment.Data()
		er := r.provider.SaveIp(ip, data)
		if err == nil && er != nil {
			err = fmt.Errorf("error saving punishments: %w", er)
		}
	}
	return err
}

// Close ...
func (r *Registry) Close() error {
	return r.Save()
}
