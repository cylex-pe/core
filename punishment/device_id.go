package punishment

import (
	"golang.org/x/exp/slices"
	"sync"
)

// Device holds all device related punishments for a user as well as their aliases.
type Device struct {
	// device is the specific identifier for this deviceid
	device string
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

func NewDevice(device string, aliases []Alias, currBan Punishment, pastBans []Punishment, currMute Punishment, pastMutes []Punishment) Device {
	return Device{
		device:      device,
		aliases:     aliases,
		currentBan:  currBan,
		pastBans:    pastBans,
		currentMute: currMute,
		pastMutes:   pastMutes,
	}
}

// Identifier ...
func (d *Device) Identifier() any {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.device
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
func (d *Device) BanHistory() []Punishment {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.pastBans
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
func (d *Device) MuteHistory() []Punishment {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.pastMutes
}
