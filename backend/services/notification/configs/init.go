package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App    `mapstructure:"app"`
	Client `mapstructure:"client"`
	Broker `mapstructure:"broker"`
	Mailer `mapstructure:"mailer"`
	Tracer `mapstructure:"tracer"`
}

type App struct {
	Name string `mapstructure:"name"`
	Env  string
}

type Client struct {
	BaseURL string `mapstructure:"base_url"`
}

type Broker struct {
	Brokers     string `mapstructure:"brokers"`
	MaxBytes    int    `mapstructure:"max_bytes"`
	MaxAttempts int    `mapstructure:"max_attempts"`
	BaseDelay   int    `mapstructure:"base_delay"`
}

type Mailer struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	From string `mapstructure:"from"`
}

type Tracer struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Endpoint string
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
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	cfg.App.Env = env
	cfg.Tracer.Endpoint = fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port)

	return &cfg, nil
}
