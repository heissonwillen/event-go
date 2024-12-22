package config

import "os"

type Config struct {
	ListenAddr        string
	SQLiteDBPath      string
	BasicAuthUser     string
	BasicAuthPassword string
}

func LoadConfig() Config {
	config := Config{
		ListenAddr:        os.Getenv("LISTEN_ADDR"),
		SQLiteDBPath:      os.Getenv("SQLITE_DB_PATH"),
		BasicAuthUser:     os.Getenv("BASIC_AUTH_USER"),
		BasicAuthPassword: os.Getenv("BASIC_AUTH_PASSWORD"),
	}

	if config.ListenAddr == "" {
		config.ListenAddr = ":8080"
	}

	if config.SQLiteDBPath == "" {
		config.SQLiteDBPath = "event-go.sqlite"
	}

	return config
}
