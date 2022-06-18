package sql

import (
	"fmt"
)

// Config provides a database configuration.
// Not all fields will be used dependent on database dialect being used.
type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
	Dialect  string `json:"dialect"`
}

// DSN returns a dsn string for the selected database dialect specified by the config.
// Supports mysql and postgres syntax but others can be added.
func (cfg Config) DSN() string {
	switch cfg.Dialect {
	case "mysql":
		// user:password@tcp(host:port)/dbname.
		// parseTime to convert to time.Time instead of []byte.
		return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true&charset=utf8mb4,utf8&multiStatements=true",
			cfg.User, cfg.Password, cfg.Protocol, cfg.Host, cfg.Port, cfg.DBName)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	default:
		return ""
	}
}

func DefaultConfig() Config {
	return Config{
		User:     "root",
		Password: "",
		Protocol: "tcp",
		Host:     "localhost",
		Port:     3306,
		DBName:   "goafweb",
		Dialect:  "mysql",
	}
}
