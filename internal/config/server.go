package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
}

func ConfigServerInit() {
	var cfg ServerConfig
	var cfgFromFlags ServerConfig

	log.Println("Server started")

	cfg.StoreInterval = time.Second * 300
	// load environment variables

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server Config environment:%+v", cfg)

	// load flags
	flag.StringVar(&cfgFromFlags.Address, "a", inst.ServerAddress, "address to listen on")
	flag.StringVar(&cfgFromFlags.Key, "k", "", "key for hash calculate")
	flag.DurationVar(&cfgFromFlags.StoreInterval, "i", inst.StoreInterval, "store interval")
	flag.StringVar(&cfgFromFlags.StoreFile, "f", inst.StoreFile, "store file")
	flag.BoolVar(&cfgFromFlags.Restore, "r", inst.Restore, "restore from file on start")
	flag.StringVar(&cfgFromFlags.DatabaseDSN, "d", inst.DatabaseDSN, "Database DSN")
	flag.Parse()
	log.Printf("Server Config flags:%+v", cfgFromFlags)

	// assign work parameters
	if cfg.Address != "" {
		inst.ServerAddress = cfg.Address
	} else {
		inst.ServerAddress = cfgFromFlags.Address
	}
	log.Printf("Server Strated with variables: address=%v", inst.ServerAddress)

	if cfg.StoreInterval >= 0 {
		inst.StoreInterval = cfg.StoreInterval
	} else {
		inst.StoreInterval = cfgFromFlags.StoreInterval
	}
	log.Printf("Server Strated with variables: StoreInterval=%v", inst.StoreInterval)

	if len(cfg.StoreFile) > 0 {
		inst.StoreFile = cfg.StoreFile
	} else {
		inst.StoreFile = cfgFromFlags.StoreFile
	}
	log.Printf("Server Strated with variables: StoreFile=%v", inst.StoreFile)

	if cfg.Key != "" {
		inst.Key = cfg.Key
	} else {
		inst.Key = cfgFromFlags.Key
	}
	log.Printf("Server Strated with variables: Key=%v", inst.Key)

	if cfg.Restore {
		inst.Restore = true
	} else if cfgFromFlags.Restore {
		inst.Restore = true
	} else {
		inst.Restore = false
	}

	if cfg.DatabaseDSN != "" {
		inst.DatabaseDSN = cfg.DatabaseDSN
	} else if cfgFromFlags.DatabaseDSN != "" {
		inst.DatabaseDSN = cfgFromFlags.DatabaseDSN
	}

	log.Printf("Server Strated with variables: Restore=%v", inst.Restore)
	log.Printf("Server Strated with variables: DatabaseDSN=%v", inst.DatabaseDSN)

}
