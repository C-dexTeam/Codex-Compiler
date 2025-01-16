package models

type Application struct {
	Site           string `mapstructure:"site"`
	Mode           string `mapstructure:"mode"`
	DevMode        bool   `mapstructure:"devMode"`
	MigrationsPath string `mapstructure:"migrationsPath"`
}
