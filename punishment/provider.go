package punishment

// Provider represents a data provider for punishments. Punishment data will be loaded and saved using this provider.
type Provider interface {
	// LoadXbox is called to retrieve a specific Users punishment by xuid. If punishments don't exist
	// it should return a Data struct with no punishments.
	Load(ptype string, identifier string) (*Container, error)
	// SaveXbox is called when saving a users xbox punishment.
	Save(ptype string, identifier string, data interface{}) error
	// LoadIp is called to retrieve a specific ip punishment registry. If punishments don't exist it should return
	// a data struct with no punishments.
	LoadIp(ip string) (IpData, error)
	// SaveIp is called when saving an instance of an ip punishment holder.
	SaveIp(ip string, ipData IpData) error
}
