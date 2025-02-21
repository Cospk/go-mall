package config

import "time"

type Config struct {
	Name     string `mapstructure:"name" json:"name"`
	App      *appConfig
	Database *databaseConfig
}

type appConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type databaseConfig struct {
	Type        string        `mapstructure:"type"`
	DSN         string        `mapstructure:"dsn"`
	MaxOpenConn int           `mapstructure:"maxopen""`
	MaxIdleConn int           `mapstructure:"maxidle"`
	MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
}
