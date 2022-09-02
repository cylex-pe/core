package punishment

// Provider represents a data provider for punishments. Punishment data will be loaded and saved using this provider.
type Provider interface {
	// LoadXbox is called to retrieve a specific Users punishment by xuid. If punishments don't exist
	// it should return a Dataer implementor with no punishments.
	Load(ptype string, identifier any) (Container, error)
	// SaveXbox is called when saving a users xbox punishment.
	Save(ptype string, identifier any, data Dataer) error
}
