package configs

import (
	"time"
)

type Config struct {
	Database struct {
		Dsn                   string        `yaml:"dsn"`
		MaxOpenConnections    int           `yaml:"maxOpenConnections"`
		MaxIdleConnections    int           `yaml:"maxIdleConnections"`
		ConnectionMaxIdleTime time.Duration `yaml:"connectionMaxIdleTime"`
		ConnectionMaxLifeTime time.Duration `yaml:"connectionMaxLifeTime"`
	} `yaml:"database"`

	Service struct {
		Port    int    `yaml:"port"`
		LogPath string `yaml:"logPath"`
	} `yaml:"service"`

	Frontend struct {
		Port int `yaml:"port"`
	}

	Nats struct {
		Url string `yaml:"url"`
	}
}

func (c Config) GetDSN() string {
	return c.Database.Dsn
}

func (c Config) GetMaxOpenConn() int {
	return c.Database.MaxOpenConnections
}

func (c Config) GetMaxIdleConn() int {
	return c.Database.MaxIdleConnections
}

func (c Config) GetConnMaxLifetime() time.Duration {
	return c.Database.ConnectionMaxIdleTime
}

func (c Config) GetConnMaxIdleTime() time.Duration {
	return c.Database.ConnectionMaxLifeTime
}

func (c Config) GetServerPort() int {
	return c.Service.Port
}

func (c Config) GetFrontendPort() int {
	return c.Frontend.Port
}

func (c Config) GetNATSUrl() string {
	return c.Nats.Url
}

func (c Config) GetLogPath() string {
	return c.Service.LogPath
}
