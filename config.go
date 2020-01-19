package main

// Config is the struct that holds the lobby server configuration
type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	Relay map[string]string
	Blacklist BlacklistConfig
}

// ServerConfig holds the basic server config.
type ServerConfig struct {
	Address string
	GeoLite2Path string
}

// DatabaseConfig holds the database config.
type DatabaseConfig struct {
	Type string
	Connection string
}

// BlacklistConfig configures the different blacklists.
type BlacklistConfig struct {
	Strings []string // General blacklisted words as RE
	IPs []string
}
