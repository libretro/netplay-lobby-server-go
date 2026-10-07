package main

// Config is the struct that holds the lobby server configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Relay     map[string]string
	Blacklist BlacklistConfig
	CORS      CORSConfig
}

// ServerConfig holds the basic server config.
type ServerConfig struct {
	Address      string
	GeoLite2Path string
	TemplatePath string
}

// DatabaseConfig holds the database config.
type DatabaseConfig struct {
	Type       string
	Connection string
}

// BlacklistConfig configures the different blacklists.
type BlacklistConfig struct {
	Strings []string // General blacklisted words as RE
	IPs     []string
}

// CORSConfig configures which origins can use /list and /add from a browser.
type CORSConfig struct {
	Origins []string // Defaults to all origins
}
