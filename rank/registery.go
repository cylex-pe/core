package rank

// Registry holds global registry of ranks. Ranks are identified by their names so each name should be unique.
type Registry struct {
	ranks map[string]Rank
}

// New returns a new Registry of ranks.
func New(ranks ...Rank) Registry {
	r := Registry{ranks: map[string]Rank{}}
	for _, rank := range ranks {
		r.ranks[rank.Name()] = rank
	}
	return r
}

// Rank returns a rank by a specific name if it has been registered.
func (r *Registry) Rank(name string) (Rank, bool) {
	if _, ok := r.ranks[name]; ok {
		return r.ranks[name], true
	}
	return nil, false
}
