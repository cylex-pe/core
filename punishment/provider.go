package punishment

type Provider interface {
	// LoadXbox is called to retrieve a specific Users punishment by xuid. If punishments don't exist
	// it should return a Data struct with no punishments.
	LoadXbox(xuid string) (XboxData, error)
	// SaveXbox is called when saving a users xbox punishment.
	SaveXbox(xuid string, xbox XboxData) error
}
