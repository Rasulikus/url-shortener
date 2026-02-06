package config

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Storage string

const (
	StoragePostgres Storage = "postgresql"
	StorageMemory   Storage = "memory"
)

func ParseStorage(s string) (Storage, error) {
	v := Storage(strings.ToLower(strings.TrimSpace(s)))
	switch v {
	case StoragePostgres, StorageMemory:
		return v, nil
	default:
		return "", fmt.Errorf("unknown storage: %q", s)
	}
}

const (
	keyLogLevel = "LOG_LEVEL"

	keyBaseURL = "BASE_URL"

	keyStorage = "STORAGE"

	keyHTTPHost         = "HTTP_HOST"
	keyHTTPPort         = "HTTP_PORT"
	keyHTTPReadTimeout  = "HTTP_READ_TIMEOUT"
	keyHTTPWriteTimeout = "HTTP_WRITE_TIMEOUT"
	keyHTTPIdleTimeout  = "HTTP_IDLE_TIMEOUT"

	keyDBHost    = "DB_HOST"
	keyDBPort    = "DB_PORT"
	keyDBUser    = "DB_USER"
	keyDBPass    = "DB_PASS"
	keyDBName    = "DB_NAME"
	keyDBSSLMode = "DB_SSLMODE"

	keyPGMinConns        = "PG_MIN_CONNS"
	keyPGMaxConns        = "PG_MAX_CONNS"
	keyPGMaxConnLifeTime = "PG_MAX_CONN_LIFETIME"
	keyPGMaxIdleTime     = "PG_MAX_CONN_IDLE_TIME"
	keyShutdownTimeout   = "SHUTDOWN_TIMEOUT"
)

type HTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DBConfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string

	MinConns        int
	MaxConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func (cfg DBConfig) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Pass),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   path.Join("/", cfg.Name),
	}
	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

type Config struct {
	LogLevel string
	BaseURL  string
	Storage  Storage // memory|postgresql

	HTTP HTTPConfig
	DB   *DBConfig

	ShutdownTimeout time.Duration
}

func getEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return value, nil
}

func getEnvInt(key string) (int, error) {
	value, err := getEnv(key)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not an integer: %q: %w", key, value, err)
	}
	return num, nil
}

func getEnvDuration(key string) (time.Duration, error) {
	value, err := getEnv(key)
	if err != nil {
		return 0, err
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid duration: %s: %w", key, value, err)
	}
	return d, nil
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	cfg := new(Config)
	var err error

	cfg.LogLevel, err = getEnv(keyLogLevel)
	if err != nil {
		return nil, err
	}

	cfg.BaseURL, err = getEnv(keyBaseURL)
	if err != nil {
		return nil, err
	}

	st, err := getEnv(keyStorage)
	if err != nil {
		return nil, err
	}
	cfg.Storage, err = ParseStorage(st)
	if err != nil {
		return nil, err
	}

	cfg.HTTP.Host, err = getEnv(keyHTTPHost)
	if err != nil {
		return nil, err
	}
	cfg.HTTP.Port, err = getEnv(keyHTTPPort)
	if err != nil {
		return nil, err
	}
	cfg.HTTP.ReadTimeout, err = getEnvDuration(keyHTTPReadTimeout)
	if err != nil {
		return nil, err
	}
	cfg.HTTP.WriteTimeout, err = getEnvDuration(keyHTTPWriteTimeout)
	if err != nil {
		return nil, err
	}
	cfg.HTTP.IdleTimeout, err = getEnvDuration(keyHTTPIdleTimeout)
	if err != nil {
		return nil, err
	}

	switch cfg.Storage {
	case StorageMemory:
	case StoragePostgres:
		cfg.DB = new(DBConfig)
		cfg.DB.Host, err = getEnv(keyDBHost)
		if err != nil {
			return nil, err
		}
		cfg.DB.Port, err = getEnv(keyDBPort)
		if err != nil {
			return nil, err
		}
		cfg.DB.User, err = getEnv(keyDBUser)
		if err != nil {
			return nil, err
		}
		cfg.DB.Pass, err = getEnv(keyDBPass)
		if err != nil {
			return nil, err
		}
		cfg.DB.Name, err = getEnv(keyDBName)
		if err != nil {
			return nil, err
		}
		cfg.DB.SSLMode, err = getEnv(keyDBSSLMode)
		if err != nil {
			return nil, err
		}

		cfg.DB.MinConns, err = getEnvInt(keyPGMinConns)
		if err != nil {
			return nil, err
		}
		cfg.DB.MaxConns, err = getEnvInt(keyPGMaxConns)
		if err != nil {
			return nil, err
		}
		cfg.DB.MaxConnLifetime, err = getEnvDuration(keyPGMaxConnLifeTime)
		if err != nil {
			return nil, err
		}
		cfg.DB.MaxConnIdleTime, err = getEnvDuration(keyPGMaxIdleTime)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown storage type, expected memory or postgresql: %s", cfg.Storage)
	}

	return cfg, nil
}
