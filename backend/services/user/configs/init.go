package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      `mapstructure:"app"`
	Server   `mapstructure:"server"`
	Database `mapstructure:"database"`
	Broker   `mapstructure:"broker"`
	Tracer   `mapstructure:"tracer"`
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

type Database struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Pass            string        `mapstructure:"pass"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxConns        int           `mapstructure:"max_conns"`
	MinConns        int           `mapstructure:"min_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
	DSN             string
}

type Broker struct {
	Brokers     string `mapstructure:"brokers"`
	MaxBytes    int    `mapstructure:"max_bytes"`
	MaxAttempts int    `mapstructure:"max_attempts"`
	BaseDelay   int    `mapstructure:"base_delay"`
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
	cfg.Database.DSN = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	cfg.Tracer.Endpoint = fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port)

	return &cfg, nil
}
