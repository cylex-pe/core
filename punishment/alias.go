package punishment

// Alias is an alias for a user, this holds a users name and xuid.
type Alias struct {
	Username string `json:"username"`
	Xuid     string `json:"xuid"`
	//TODO: device id
}
