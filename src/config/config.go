package config

import "os"

const (
	Production  string = "prod"
	Development string = "dev"
)

type Config struct {
	Profile         string
	PostgresConfig  PostgresConfig
	BasicAuthConfig BasicAuthConfig
}

type BasicAuthConfig struct {
	Username string
	Password string
}

func NewConfig() *Config {
	return &Config{
		Profile:         getProfile(),
		PostgresConfig:  getPostgresConfig(),
		BasicAuthConfig: getBasicAuthConfig(),
	}
}

func getBasicAuthConfig() BasicAuthConfig {
	return BasicAuthConfig{
		Username: mustGetenv("BASIC_AUTH_USER"),
		Password: mustGetenv("BASIC_AUTH_PASSWORD"),
	}
}

type PostgresConfig struct {
	Username string
	Password string
	Database string
	Hostname string
}

func (p *PostgresConfig) GetConnectionURL() string {
	return "postgres://" + p.Username + ":" + p.Password + "@" + p.Hostname + "/" + p.Database + "?sslmode=disable"
}

func getPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Username: mustGetenv("POSTGRES_USER"),
		Password: mustGetenv("POSTGRES_PASSWORD"),
		Database: mustGetenv("POSTGRES_DB"),
		Hostname: mustGetenv("POSTGRES_HOST"),
	}
}

func getProfile() string {
	switch os.Getenv("PROFILE") {
	case Development:
		return Development
	default:
		return Production
	}
}

func mustGetenv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("no env for: " + key)
	}

	return value
}
