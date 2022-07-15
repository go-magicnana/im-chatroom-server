package rediss

import "time"

type Config struct {
	// Network "tcp"
	Network string
	// Addr "127.0.0.1:6379"
	Addr string
	// Password string .If no password then no 'AUTH'. Default ""
	Password string
	// If Database is empty "" then no 'SELECT'. Default ""
	Database string
	// MaxIdle 0 no limit
	MaxIdle int
	// MaxActive 0 no limit
	MaxActive int
	// IdleTimeout  time.Duration(5) * time.Minute
	IdleTimeout time.Duration
	// Prefix "myprefix-for-this-website". Default ""
	Prefix string
}

func DefaultConfig() Config {
	return Config{
		Network:  "tcp",
		Addr:     "47.95.148.121:6379",
		Password: "o1trUmeh",
		Database: "0",
	}
}
