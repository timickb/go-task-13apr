package config

type Config struct {
	AppPort   int    `yaml:"app_port"`
	StoreFile string `yaml:"users_file"`
}

func NewDefault() *Config {
	return &Config{
		AppPort:   8080,
		StoreFile: "users.json",
	}
}
