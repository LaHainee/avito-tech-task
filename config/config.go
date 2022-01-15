package config

type ServerConfig struct {
	DatabaseConnString string `toml:"database_conn_string"`
}

type Config struct {
	LoggingLevel    string       `toml:"logging_level"`
	LoggingFilePath string       `toml:"logging_file_path"`
	CurrencyApiURL  string       `toml:"currency_api_url"`
	Server          ServerConfig `toml:"server"`
}

func NewConfig() *Config {
	return &Config{}
}
