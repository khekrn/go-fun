package config

import (
	"errors"
	"fmt"

	"gopkg.in/ini.v1"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Schema   string
	SSLMode  string
	Timezone string
}

func (d *DatabaseConfig) ConnectionURL() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s search_path=%s port=%s sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.Name, d.Schema, d.Port, d.SSLMode, d.Timezone)
}

type ServerConfig struct {
	Port     string
	Mode     string
	LogLevel string
}

func Load(fileName string) (*Config, error) {
	cfg, err := ini.Load(fileName)
	if err != nil {
		return nil, errors.New("failed to load config file: " + err.Error())
	}

	config := &Config{}

	serverSection := cfg.Section("server")
	config.Server = ServerConfig{
		Port:     serverSection.Key("port").MustString("8080"),
		Mode:     serverSection.Key("mode").MustString("debug"),
		LogLevel: serverSection.Key("logLevel").MustString("info"),
	}

	dbSection := cfg.Section("database")
	config.Database = DatabaseConfig{
		Host:     dbSection.Key("host").MustString("localhost"),
		Port:     dbSection.Key("port").MustString("5432"),
		User:     dbSection.Key("user").MustString("postgres"),
		Password: dbSection.Key("password").MustString(""),
		Name:     dbSection.Key("name").MustString("proddb"),
		Schema:   dbSection.Key("schema").MustString("public"),
		SSLMode:  dbSection.Key("sslmode").MustString("disable"),
		Timezone: dbSection.Key("timezone").MustString("UTC"),
	}

	return config, nil
}
