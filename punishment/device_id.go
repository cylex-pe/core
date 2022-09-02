package punishment

import (
	"golang.org/x/exp/slices"
	"sync"
)

// Device holds all device related punishments for a user as well as their aliases.
type Device struct {
	// aliases represent the information of other accounts that have the same device-id as this one.
	aliases []Alias

	// currentBan holds a users current device ban, if they're  not currently device banned it will be the default value.
	currentBan Punishment
	//pastBans holds all a users past device bans.
	pastBans []Punishment
	// currentMute is the players current active mute, this value will be the default value for Punishment if their not muted.
	currentMute Punishment
	// pastMutes store a history of all the users past mutes.
	pastMutes []Punishment

	lock sync.RWMutex
}

// DeviceData is a data representation of device used for loading and saving devices.
type DeviceData struct {
	// Aliases represnt aliases within Device.
	Aliases []Alias `json:"aliases"`
	// CurrentBan represents currentBan within Device.
	CurrentBan Punishment `json:"current_ban"`
	// PastBans represents pastBans within Device.
	PastBans []Punishment `json:"past_bans"`
	// CurrentMute represents currentMute within Device.
	CurrentMute Punishment `json:"current_mute"`
	// PastMutes represents pastMutes within Device.
	PastMutes []Punishment `json:"past_mutes"`
}

func (d *DeviceData) Container() Container {
	return &Ip{
		aliases:     d.Aliases,
		currentBan:  d.CurrentBan,
		pastBans:    d.PastBans,
		currentMute: d.CurrentMute,
		pastMutes:   d.PastMutes,
	}
}

// AddAlias attempts to add an alias into IP, it will return true if it managed to add it and false if a value already
// existed.
func (d *Device) AddAlias(alias Alias) bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	if !slices.Contains(d.aliases, alias) {
		d.aliases = append(d.aliases, alias)
		return true
	}
	return false
}

// Aliases returns all the associated aliases with this account.
func (d *Device) Aliases() []Alias {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.aliases
}

// Banned returns whether the current holder is banned or not.
func (d *Device) Banned() bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.currentBan != Punishment{}
}

// CurrentBan returns the users current ban.
func (d *Device) CurrentBan() Punishment {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.currentBan
}

// Ban adds a ban to a users current ban, it moves their previous current ban to their pastBans.
func (d *Device) Ban(b Punishment) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.pastBans = append(d.pastBans, d.currentBan)
	d.currentBan = b
}

// BanHistory returns the BanHistory for the user, instead of returning a ban it returns a []BanData as it's meant to
// be used for reading bans only.
func (d *Device) BanHistory() []Data {
	var bans []Data
	for _, ban := range d.pastBans {
		bans = append(bans, ban.Data())
	}
	return bans
}

// Muted returns whether the current holder is banned or not.
func (d *Device) Muted() bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.currentBan != Punishment{}
}

// CurrentMute returns the users current ban.
func (d *Device) CurrentMute() Punishment {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.currentBan
}

// Mute adds a ban to a users current ban, it moves their previous current ban to their pastBans.
func (d *Device) Mute(b Punishment) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.pastBans = append(d.pastBans, d.currentBan)
	d.currentBan = b
}

// MuteHistory returns the MuteHistory for the user, instead of returning a punishment it returns a []Data as it's meant to
// be used for reading mutes only.
func (d *Device) MuteHistory() []Data {
	var bans []Data
	for _, ban := range d.pastBans {
		bans = append(bans, ban.Data())
	}
	return bans
}

// Data returns the data representation of IP.
func (d *Device) Data() Dataer {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return &IpData{
		Aliases:     d.aliases,
		CurrentBan:  d.currentBan,
		PastBans:    d.pastBans,
		CurrentMute: d.currentMute,
		PastMutes:   d.pastMutes,
	}
}
