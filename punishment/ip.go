package punishment

import (
	"golang.org/x/exp/slices"
	"sync"
)

// Ip holds all ip related punishments for a user as well as their aliases.
type Ip struct {
	// aliases represent the information of other accounts that have the same ip address as this one.
	aliases []Alias

	// currentBan holds a users current Ip ban, if their not currently IPBanned it will be the default value.
	currentBan Punishment
	//pastBans holds all a users past Ip bans.
	pastBans []Punishment

	lock sync.RWMutex
}

// IpData is a data representation of Ip used for loading and saving Ip.
type IpData struct {
	Aliases    []Alias      `json:"aliases"`
	CurrentBan Punishment   `json:"current_ban"`
	PastBans   []Punishment `json:"past_bans"`
}

// IpAlias represents an alias with respect to a users Ip address.
type IpAlias struct {
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
func (i *Ip) BanHistory() []Data {
	var bans []Data
	for _, ban := range i.pastBans {
		bans = append(bans, ban.Data())
	}
	return bans
}

// Data returns the data representation of IP.
func (i *Ip) Data() IpData {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return IpData{
		Aliases:    i.aliases,
		CurrentBan: i.currentBan,
		PastBans:   i.pastBans,
	}
}
