package internal

import (
	"cmp"
	"flag"
	"os"
)

type Config struct {
	Addr  string
	Port  int
	DBDSN string
}

const (
	defaultAddr  = "localhost"
	defaulPort   = 8080
	defaultDBStr = "postgresql://user:password@localhost:5432/library?sslmode=disable"
)

func ReadConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.Addr, "addr", defaultAddr, "flag for use custom server addr")
	flag.IntVar(&cfg.Port, "port", defaulPort, "flag for use custom server port")
	flag.StringVar(&cfg.DBDSN, "db", defaultDBStr, "flag for setup dv connection string")
	flag.Parse()

	cfg.DBDSN = cmp.Or(os.Getenv("DB_DSN"), cfg.DBDSN)
	return &cfg
}
