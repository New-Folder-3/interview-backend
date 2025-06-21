package conf

import (
	"interview-backend/cmd/flags"
	"path/filepath"
)

type Database struct {
	Type        string `json:"type" env:"TYPE"`
	Host        string `json:"host" env:"HOST"`
	Port        int    `json:"port" env:"PORT"`
	User        string `json:"user" env:"USER"`
	Password    string `json:"password" env:"PASSWORD"`
	Name        string `json:"name" env:"NAME"`
	TablePrefix string `json:"tablePrefix" env:"TABLE_PREFIX"`
	SSLMode     string `json:"sslMode" env:"SSL_MODE"`
	DSN         string `json:"dsn" env:"DSN"`
	DBFile      string `json:"db_file" env:"DB_FILE"`
}

type Schema struct {
	Port   int    `json:"port" env:"PORT"`
	Listen string `json:"listen" env:"LISTEN"`
}

type Config struct {
	Database Database `json:"database" envPrefix:"DB_"`
	Schema   Schema   `json:"schema" envPrefix:"SCHEMA"`
}

func DefaultConfig() *Config {
	dbPath := filepath.Join(flags.DataDir, "data.db")
	return &Config{
		Database: Database{
			Type:        "sqlite3",
			Port:        0,
			TablePrefix: "x_",
			DBFile:      dbPath,
		},
		Schema: Schema{
			Port:   2233,
			Listen: "127.0.0.1",
		},
	}
}
