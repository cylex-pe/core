package rank

// Rank represents a rank that a player can have.
type Rank interface {
	// Name of the rank.
	Name() string
	// Level is the importance level / security clearance of the rank.
	Level() int
	// Staff is whether this rank is a staff rank or not.
	Staff() bool
	// ChatFormat returns the chat format for the specific rank.
	ChatFormat() string
}
