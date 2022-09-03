package punishment

import "sync"

// Xbox represents a specific user's punishment information such as a users bans, ipbans, reports.
type Xbox struct {
	// currentBan is the players current active ban, this value can be null so we use a pointer.
	currentBan Punishment
	// pastBans store a history of all the users past bans.
	pastBans []Punishment
	// currentMute is the players current active mute, this value can be null so we use a pointer.
	currentMute Punishment
	// pastMutes store a history of all the users past mutes.
	pastMutes []Punishment

	lock sync.RWMutex
}

// XboxData holds the required data for an Xbox.
type XboxData struct {
	// CurrentBan represents currentBan within Xbox.
	CurrentBan Punishment `json:"current_ban"`
	// PastBans represents pastBans within Xbox.
	PastBans []Punishment `json:"past_bans"`
	// CurrentMute represents currentMute within Xbox.
	CurrentMute Punishment `json:"current_mute"`
	// PastMutes represents pastMutes within Xbox.
	PastMutes []Punishment `json:"past_mutes"`
}

func (x XboxData) Container() Container {
	return &Xbox{
		currentBan:  x.CurrentBan,
		pastBans:    x.PastBans,
		currentMute: x.CurrentMute,
		pastMutes:   x.PastMutes,
	}
}

// Banned returns whether the current holder is banned or not.
func (x *Xbox) Banned() bool {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.currentBan != Punishment{}
}

// CurrentBan returns the users current ban.
func (x *Xbox) CurrentBan() Punishment {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.currentBan
}

// Ban adds a ban to a users current ban, it moves their previous current ban to their pastBans.
func (x *Xbox) Ban(b Punishment) {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.pastBans = append(x.pastBans, x.currentBan)
	x.currentBan = b
}

// BanHistory returns the BanHistory for the user, instead of returning a ban it returns a []BanData as it's meant to
// be used for reading bans only.
func (x *Xbox) BanHistory() []Punishment {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.pastBans
}

// Muted returns whether the current holder is banned or not.
func (x *Xbox) Muted() bool {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.currentBan != Punishment{}
}

// CurrentMute returns the users current ban.
func (x *Xbox) CurrentMute() Punishment {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.currentBan
}

// Mute adds a ban to a users current ban, it moves their previous current ban to their pastBans.
func (x *Xbox) Mute(b Punishment) {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.pastBans = append(x.pastBans, x.currentBan)
	x.currentBan = b
}

// MuteHistory returns the MuteHistory for the user, instead of returning a punishment it returns a []Data as it's meant to
// be used for reading mutes only.
func (x *Xbox) MuteHistory() []Punishment {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.pastMutes
}

// Data returns the data representation for this punishment.
func (x *Xbox) Data() DataHolder {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return &XboxData{
		CurrentBan:  x.currentBan,
		PastBans:    x.pastBans,
		CurrentMute: x.currentMute,
		PastMutes:   x.pastMutes,
	}
}
