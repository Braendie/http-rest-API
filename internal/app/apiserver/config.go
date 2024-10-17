package apiserver

// Config holds the configuration settings for the server application.
// It includes the following fields:
// - BindAddr: the address the server will bind to, used for listening to incoming connections.
// - LogLevel: the logging level for the application, defining the verbosity of the logs.
// - DatabaseURL: the URL for connecting to the database, including the necessary credentials and connection string.
// - SessionKey: the secret key used for securing user sessions, which ensures the integrity of session data.
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
}

// NewConfig returns a new config with filled fields from the toml file.
func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		LogLevel:    "debug",
	}
}
