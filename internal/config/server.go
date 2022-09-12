package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type ServerConfig struct {
	Address            string        `env:"ADDRESS"`
	StoreInterval      time.Duration `env:"STORE_INTERVAL"`
	StoreFile          string        `env:"STORE_FILE"`
	Restore            bool          `env:"RESTORE"`
	Key                string        `env:"KEY"`
	DatabaseDSN        string        `env:"DATABASE_DSN"`
	LogHTTP            bool          `env:"Log_HTTP"`
	PrivateKeyFileName string        `env:"CRYPTO_KEY"`
}

// ConfigServerInit load config from flags and environment variables
func ConfigServerInit() {
	var cfgEnv ServerConfig
	var cfgFromFlags ServerConfig

	inst.Sugar.Infow("Server started")

	cfgEnv.StoreInterval = time.Second * 300
	// load environment variables

	err := env.Parse(&cfgEnv)
	if err != nil {
		log.Fatal(err)
	}

	inst.Sugar.Infof("Server Config environment:%+v", cfgEnv)

	// load flags
	flag.StringVar(&cfgFromFlags.Address, "a", inst.ServerAddress, "address to listen on")
	flag.StringVar(&cfgFromFlags.Key, "k", "", "key for hash calculate")
	flag.DurationVar(&cfgFromFlags.StoreInterval, "i", inst.StoreInterval, "store interval")
	flag.StringVar(&cfgFromFlags.StoreFile, "f", inst.StoreFile, "store file")
	flag.BoolVar(&cfgFromFlags.Restore, "r", inst.Restore, "restore from file on start")
	flag.StringVar(&cfgFromFlags.DatabaseDSN, "d", inst.DatabaseDSN, "Database DSN")
	flag.BoolVar(&cfgFromFlags.LogHTTP, "l", inst.LogHTTP, "log HTTP switch, so mach information, switch ON only for debug")
	flag.StringVar(&cfgFromFlags.PrivateKeyFileName, "crypto-key", inst.PrivateKeyFileName, "private key file name")

	flag.Parse()

	inst.Sugar.Infof("Server Config flags:%+v", cfgFromFlags)

	// assign work parameters
	if cfgEnv.Address != "" {
		inst.ServerAddress = cfgEnv.Address
	} else {
		inst.ServerAddress = cfgFromFlags.Address
	}
	inst.Sugar.Infof("Server Strated with variables: address=%v", inst.ServerAddress)

	if cfgEnv.StoreInterval >= 0 {
		inst.StoreInterval = cfgEnv.StoreInterval
	} else {
		inst.StoreInterval = cfgFromFlags.StoreInterval
	}
	inst.Sugar.Infof("Server Strated with variables: StoreInterval=%v", inst.StoreInterval)

	if len(cfgEnv.StoreFile) > 0 {
		inst.StoreFile = cfgEnv.StoreFile
	} else {
		inst.StoreFile = cfgFromFlags.StoreFile
	}
	inst.Sugar.Infof("Server Strated with variables: StoreFile=%v", inst.StoreFile)

	if cfgEnv.Key != "" {
		inst.Key = cfgEnv.Key
	} else {
		inst.Key = cfgFromFlags.Key
	}
	inst.Sugar.Infof("Server Strated with variables: Key=%v", inst.Key)

	if cfgEnv.Restore {
		inst.Restore = true
	} else if cfgFromFlags.Restore {
		inst.Restore = true
	} else {
		inst.Restore = false
	}
	inst.Sugar.Infof("Server Strated with variables: Restore=%v", inst.Restore)

	if cfgEnv.DatabaseDSN != "" {
		inst.DatabaseDSN = cfgEnv.DatabaseDSN
	} else if cfgFromFlags.DatabaseDSN != "" {
		inst.DatabaseDSN = cfgFromFlags.DatabaseDSN
	}
	inst.Sugar.Infof("Server Strated with variables: DatabaseDSN=%v", inst.DatabaseDSN)

	if cfgEnv.LogHTTP {
		inst.LogHTTP = true
	} else if cfgFromFlags.LogHTTP {
		inst.LogHTTP = true
	} else {
		inst.LogHTTP = false
	}
	inst.Sugar.Infof("Server Strated with variables: LogHTTP=%v", inst.LogHTTP)

	if cfgEnv.PrivateKeyFileName != "" {
		inst.PrivateKeyFileName = cfgEnv.PrivateKeyFileName
	} else if cfgFromFlags.PrivateKeyFileName != "" {
		inst.PrivateKeyFileName = cfgFromFlags.PrivateKeyFileName
	}
	inst.Sugar.Infof("Server Strated with variables: PrivateKeyFileName=%v", inst.PrivateKeyFileName)
}
