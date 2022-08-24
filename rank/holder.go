package rank

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Holder holds a list of ranks a player has.
type Holder struct {
	// ranks is the list of ranks the player has.
	ranks map[string]Rank
	// permission represents the highest permission a Rank Holder has.
	permission int
	// staff represents whether this user is a staff member or not.
	staff bool
	// lock locks the data for accessing.
	lock sync.RWMutex
}

// NewHolder returns a new rank holder.
func NewHolder() Holder {
	return Holder{
		ranks:      map[string]Rank{},
		permission: 0,
		staff:      false,
	}
}

// Add adds a rank to the user.
func (h *Holder) Add(rank Rank) {
	defer h.lock.Unlock()
	h.lock.Lock()
	h.ranks[rank.Name()] = rank
	h.recalculatePermission()
}

// Remove removes a rank from the holder.
func (h *Holder) Remove(name string) bool {
	defer h.lock.Unlock()
	h.lock.Lock()
	if _, ok := h.ranks[name]; ok {
		delete(h.ranks, name)
		return true
	}
	return false
}

// Rank returns a users rank, returning nil if they don't have one.
func (h *Holder) Rank(name string) *Rank {
	defer h.lock.RUnlock()
	h.lock.RLock()
	if _, ok := h.ranks[name]; ok {
		rank := h.ranks[name]
		return &rank
	}
	return nil
}

// Staff returns whether this rank holder is a staff member or not.
func (h *Holder) Staff() bool {
	defer h.lock.RUnlock()
	h.lock.RLock()
	return h.staff
}

// HolderFromJson returns a new holder from Marshaled data, I was going do this in Unmarshal, but I need access to a
// register.
func HolderFromJson(data []byte, registry *Registry) (*Holder, error) {
	aux := &struct {
		Ranks []string `json:"ranks"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return nil, err
	}
	holder := NewHolder()
	for _, n := range aux.Ranks {
		rank, ok := registry.Rank(n)
		if !ok {
			return nil, fmt.Errorf("Unregistered Rank")
		}
		holder.ranks[rank.Name()] = rank
		if rank.Level() > holder.permission {
			holder.permission = rank.Level()
		}
		if rank.Staff() && !holder.staff {
			holder.staff = true
		}
	}
	return &holder, nil
}

// recalculatePermission recalculates a Holders the highest permission, callers of this method should have the mutex within
// Holder locked.
func (h *Holder) recalculatePermission() {
	h.permission = 0
	h.staff = false
	for _, rank := range h.ranks {
		if rank.Level() > h.permission {
			h.permission = rank.Level()
		}
		if rank.Staff() && !h.staff {
			h.staff = true
		}
	}
}

// MarshalJSON ...
func (h *Holder) MarshalJSON() ([]byte, error) {
	var ranks []string
	for rank := range h.ranks {
		ranks = append(ranks, rank)
	}
	return json.Marshal(&struct {
		Ranks []string `json:"ranks"`
	}{
		Ranks: ranks,
	})
}
