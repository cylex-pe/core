package punishment

import "fmt"

type Registry struct {
	provider        Provider
	xboxPunishments map[string]*Xbox
}

// New returns a new punishment handler.
func New(provider Provider) Registry {
	return Registry{
		provider: provider,
	}
}

// Xbox attempts to load an xbox xuid punishment holder from the handler or the database.
func (r *Registry) Xbox(xuid string) (*Xbox, error) {
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

// Save attempts to save all the data within the punishment Registry.
func (r *Registry) Save() error {
	var err error
	err = nil
	for xuid, punishment := range r.xboxPunishments {
		data := punishment.Data()
		er := r.provider.SaveXbox(xuid, data)
		if err == nil && er != nil {
			err = er
		}
	}
	return err
}

// Close ...
func (r *Registry) Close() error {
	return r.Save()
}
