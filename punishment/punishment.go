package punishment

import "time"

// Punishment represents a generic punishment on a player.
type Punishment struct {
	// time is the time when this punishment was issued.
	time int
	// reason is the reason for this specific punishment.
	reason string
	// issuer is the name of the specific user that issued this punishment
	issuer string
	// expires holds a bool saying whether this punishment expires or not.
	expires bool
	// expirationTime represents when this punishment expires, it's irrelevant unless Expires is true.
	expirationTime int
}

// Data holds the required data for a Punishment.
type Data struct {
	// Time represents Time within Punishment.
	Time int `json:"time"`
	// Reason represents Reason within Punishment.
	Reason string `json:"reason"`
	// Issuer  represents Banner within Punishment.
	Issuer string `json:"banner"`
	// Expires represents Expires within Punishment.
	Expires bool `json:"expires"`
	// ExpirationTime  represents expirationTime within Punishment.
	ExpirationTime int `json:"duration"`
}

// NewPunishment returns a new Ban object.
func NewPunishment(time int, reason string, issuer string) Punishment {
	return Punishment{
		time:   time,
		reason: reason,
		issuer: issuer,
	}
}

// Issuer returns the user who did this specific ban.
func (p *Punishment) Issuer() string {
	return p.issuer
}

// Reason returns the reason for the ban.
func (p *Punishment) Reason() string {
	return p.reason
}

// Notes returns the notes for the ban.
func (p *Punishment) Notes() string {
	return p.Notes()
}

// Expired tells weather a specific ban has expired or not
func (p *Punishment) Expired() bool {
	if !p.expires {
		return false
	}
	return time.Now().Second() > p.expirationTime
}

// Data returns the ban data representation
func (p *Punishment) Data() Data {
	return Data{
		Time:           p.time,
		Reason:         p.reason,
		Issuer:         p.issuer,
		Expires:        p.expires,
		ExpirationTime: p.expirationTime,
	}
}
