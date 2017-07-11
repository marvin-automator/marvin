package domain

// This struct stores... config variables. What a surprise!
type Config struct {
	// Whether users need to log in to use Marvin.
	// By default, we're assuming it's running on a local network,
	// used by a single user, and therefore doesn't need multiple accounts.
	AccountsEnabled bool
}

// DefaultConfig stores the default values for all our configuration options
var DefaultConfig = Config{
	AccountsEnabled: true,
}

// The ConfigStore interface can handle persisting a config object.
type ConfigStore interface {
	GetConfig() (Config, error)
}
