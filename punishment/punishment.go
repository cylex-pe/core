package punishment

import "time"

// Punishment represents a generic punishment on a player.
type Punishment struct {
	// Time is the Time when this punishment was issued.
	Time int `json:"time"`
	// Reason is the Reason for this specific punishment.
	PunishmentReason string `json:"reason"`
	// issuer is the name of the specific user that issued this punishment
	PunishmentIssuer string `json:"banner"`
	// expires holds a bool saying whether this punishment expires or not.
	Expires bool `json:"expires"`
	// ExpirationTime represents when this punishment expires, it's irrelevant unless Expires is true.
	ExpirationTime int `json:"duration"`
}

// NewPunishment returns a new PunishmentReason object.
func NewPunishment(time int, reason string, issuer string) Punishment {
	return Punishment{
		Time:             time,
		PunishmentReason: reason,
		PunishmentIssuer: issuer,
	}
}

// Issuer returns the user who did this specific ban.
func (p *Punishment) Issuer() string {
	return p.PunishmentIssuer
}

// Reason returns the reason for the ban.
func (p *Punishment) Reason() string {
	return p.PunishmentReason
}

// Notes returns the notes for the ban.
func (p *Punishment) Notes() string {
	return p.Notes()
}

// Expired tells weather a specific ban has expired or not
func (p *Punishment) Expired() bool {
	if !p.Expires {
		return false
	}
	return time.Now().Second() > p.ExpirationTime
}
