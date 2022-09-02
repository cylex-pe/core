package punishment

import (
	"golang.org/x/exp/slices"
	"sync"
)

// Ip holds all ip related punishments for a user as well as their aliases.
type Ip struct {
	// ip the ip address that identifies this ip.
	ip string
	// aliases represent the information of other accounts that have the same ip address as this one.
	aliases []Alias
	// currentBan holds a users current Ip ban, if their not currently IPBanned it will be the default value.
	currentBan Punishment
	//pastBans holds all a users past Ip bans.
	pastBans []Punishment
	// currentMute is the players current active mute, if they're current not ip muted it will be the default value.
	currentMute Punishment
	// pastMutes store a history of all the users past mutes.
	pastMutes []Punishment

	lock sync.RWMutex
}

func NewIp(ip string, aliases []Alias, currBan Punishment, pastBans []Punishment, currMute Punishment, pastMutes []Punishment) Ip {
	return Ip{
		ip:          ip,
		aliases:     aliases,
		currentBan:  currBan,
		pastBans:    pastBans,
		currentMute: currMute,
		pastMutes:   pastMutes,
	}
}

// Identifier ...
func (i *Ip) Identifier() any {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.ip
}

// AddAlias attempts to add an alias into IP, it will return true if it managed to add it and false if a value already
// existed.
func (i *Ip) AddAlias(alias Alias) bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	if !slices.Contains(i.aliases, alias) {
		i.aliases = append(i.aliases, alias)
		return true
	}
	return false
}

// Aliases returns all the associated aliases with this account.
func (i *Ip) Aliases() []Alias {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.aliases
}

// Banned returns whether the current holder is banned or not.
func (i *Ip) Banned() bool {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.currentBan != Punishment{}
}

// CurrentBan returns the users current ban.
func (i *Ip) CurrentBan() Punishment {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.currentBan
}

// Ban adds a ban to a users current ban, it moves their previous current ban to their pastBans.
func (i *Ip) Ban(b Punishment) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.pastBans = append(i.pastBans, i.currentBan)
	i.currentBan = b
}

// BanHistory returns the BanHistory for the user, instead of returning a ban it returns a []BanData as it's meant to
// be used for reading bans only.
func (i *Ip) BanHistory() []Punishment {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.pastBans
}

// Muted returns whether the current holder is banned or not.
func (i *Ip) Muted() bool {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.currentBan != Punishment{}
}

// CurrentMute returns the users current ban.
func (i *Ip) CurrentMute() Punishment {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.currentBan
}

// Mute adds a ban to a users current ban, it moves their previous current ban to their pastBans.
func (i *Ip) Mute(b Punishment) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.pastBans = append(i.pastBans, i.currentBan)
	i.currentBan = b
}

// MuteHistory returns the MuteHistory for the user, instead of returning a punishment it returns a []Data as it's meant to
// be used for reading mutes only.
func (i *Ip) MuteHistory() []Punishment {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.pastMutes
}
