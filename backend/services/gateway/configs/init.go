package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App     `mapstructure:"app"`
	Server  `mapstructure:"server"`
	Service `mapstructure:"service"`
}

type App struct {
	Name string `mapstructure:"name"`
	Env  string
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	Timeout struct {
		Read     time.Duration `mapstructure:"read"`
		Write    time.Duration `mapstructure:"write"`
		Shutdown time.Duration `mapstructure:"shutdown"`
	} `mapstructure:"timeout"`
}

type Service struct {
	Auth struct {
		Addr string
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"auth"`
}

func Init(path string) (*Config, error) {
	if path == "" {
		path = "./configs"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // default
	}

	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.App.Env = env
	cfg.Auth.Addr = fmt.Sprintf("%s:%d", cfg.Auth.Host, cfg.Auth.Port)

	return &cfg, nil
}
