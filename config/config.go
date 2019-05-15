package config

import (
	"errors"
	"flag"

	"github.com/spf13/viper"
)

// Config holds all application configurations
type Config struct {
	Env         string              `json:"env"`
	Application *ApplicationSetting `json:"application"`
	Database    *DatabaseSetting    `json:"database"`
	Auth        *AuthSetting        `json:"auth"`
}

// ApplicationSetting holds all general application configurations
type ApplicationSetting struct {
	Port           int `json:"port"`
	RequestTimeout int `json:"requestTimeout"`
}

// DatabaseSetting holds all database configurations
type DatabaseSetting struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// AuthSetting holds all auth related configurations
type AuthSetting struct {
	Jwt *JwtSetting `json:"jwt"`
}

// JwtSetting holds all JWT related configurations
type JwtSetting struct {
	Secret          string `json:"secret"`
	ExpiryInSeconds int    `json:"expiryInSeconds"`
}

// Load will fetch configuration from environment specific file and populate the configuration struct.
func (config *Config) Load() error {
	var env string
	flag.StringVar(&env, "env", "dev", "environment")

	flag.Parse()

	viperRegistry := viper.New()
	viperRegistry.AddConfigPath("./config")
	viperRegistry.SetConfigName(env)
	viperRegistry.SetConfigType("json")
	viperRegistry.SetEnvPrefix("todo")
	viperRegistry.AutomaticEnv()

	config.Env = env

	if err := viperRegistry.ReadInConfig(); err != nil {
		return err
	}

	if err := config.configureApplication(viperRegistry); err != nil {
		return err
	}

	if err := config.configureDB(viperRegistry); err != nil {
		return err
	}

	if err := config.configureAuth(viperRegistry); err != nil {
		return err
	}

	return nil
}

// configureDB loads database specific configuration.
// Database address, name, user or password will be taken from OS environment
// variables if its not provided in the JSON config. If case OS environment
// variables are missing, an error will be thrown.
func (config *Config) configureDB(viperRegistry *viper.Viper) error {
	dbConfig := &DatabaseSetting{}
	db := viperRegistry.Sub("db")

	if err := db.Unmarshal(dbConfig); err != nil {
		return err
	}

	if !db.IsSet("user") {
		if !viperRegistry.IsSet("DB_USER") {
			return errors.New("database user not set")
		}

		dbConfig.User = viperRegistry.GetString("DB_USER")
	}

	if !db.IsSet("password") {
		if !viperRegistry.IsSet("DB_PASSWORD") {
			return errors.New("database password not set")
		}

		dbConfig.Password = viperRegistry.GetString("DB_PASSWORD")
	}

	if !db.IsSet("name") {
		if !viperRegistry.IsSet("DB_NAME") {
			return errors.New("database name not set")
		}

		dbConfig.Name = viperRegistry.GetString("DB_NAME")
	}

	if !db.IsSet("address") {
		if !viperRegistry.IsSet("DB_ADDRESS") {
			return errors.New("database address not set")
		}

		dbConfig.Address = viperRegistry.GetString("DB_ADDRESS")
	}

	config.Database = dbConfig

	return nil
}

// configureApplication will load general application configurations.
// Port and request timeout values are optional (default values are applied).
func (config *Config) configureApplication(viperRegistry *viper.Viper) error {
	appConfig := &ApplicationSetting{}
	appSettings := viperRegistry.Sub("application")

	if appSettings == nil {
		return errors.New("application section missing in configuration")
	}

	if err := appSettings.Unmarshal(appConfig); err != nil {
		return err
	}

	if !appSettings.IsSet("port") {
		appConfig.Port = 8080
	}

	if !appSettings.IsSet("requestTimeout") {
		appConfig.RequestTimeout = 10
	}

	config.Application = appConfig

	return nil
}

// configureAuth loads auth specific configurations.
// JWT secret is taken from OS environment variable, if missing.
// JWT expiry time is optional (defaults applied).
func (config *Config) configureAuth(viperRegistry *viper.Viper) error {
	authConfig := &AuthSetting{}
	authSettings := viperRegistry.Sub("auth")

	if authSettings == nil {
		return errors.New("auth section missing in configuration")
	}

	jwtConfig := &JwtSetting{}

	jwtSettings := authSettings.Sub("jwt")

	if jwtSettings == nil {
		return errors.New("jwt section missing in configuration")
	}

	if err := jwtSettings.Unmarshal(jwtConfig); err != nil {
		return err
	}

	if !jwtSettings.IsSet("secret") {
		if !viperRegistry.IsSet("AUTH_JWT_SECRET") {
			return errors.New("jwt secret not set")
		}
	}

	if !jwtSettings.IsSet("expiryInSeconds") {
		jwtConfig.ExpiryInSeconds = 3600
	}

	authConfig.Jwt = jwtConfig
	config.Auth = authConfig

	return nil
}
