package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type ServerConfig struct {
	Address            string        `env:"ADDRESS" json:"address"`
	StoreInterval      time.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile          string        `env:"STORE_FILE" json:"store_file"`
	Restore            bool          `env:"RESTORE" json:"restore"`
	Key                string        `env:"KEY" json:"key"`
	DatabaseDSN        string        `env:"DATABASE_DSN" json:"database_dsn"`
	LogHTTP            bool          `env:"Log_HTTP" json:"log_http"`
	PrivateKeyFileName string        `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile         string        `env:"CONFIG"`
}

// ConfigServerInit load config from flags and environment variables
func ConfigServerInit() {
	var cfgEnv ServerConfig
	var cfgFromFlags ServerConfig
	var cfgFromJsonFile ServerConfig

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
	flag.StringVar(&cfgFromFlags.ConfigFile, "c", "", "config file name")

	flag.Parse()

	inst.Sugar.Infof("Server Config flags:%+v", cfgFromFlags)

	// load config file if exist
	var ConfigFileName string
	if cfgEnv.ConfigFile != "" {
		ConfigFileName = cfgEnv.ConfigFile
	} else {
		ConfigFileName = cfgFromFlags.ConfigFile
	}
	if ConfigFileName != "" {
		content, err := os.ReadFile(ConfigFileName)
		if err != nil {
			inst.Sugar.Fatal("Error when opening file: ", err)
		}

		// Now let's unmarshall the data into `payload`
		err = json.Unmarshal(content, &cfgFromJsonFile)
		if err != nil {
			inst.Sugar.Fatal("Error during Unmarshal(): ", err)
		}
		inst.Sugar.Infof("Agent Config from file:%+v", cfgFromJsonFile)
	}

	// assign work parameters
	if cfgFromJsonFile.Address != "" {
		inst.ServerAddress = cfgFromJsonFile.Address
	}
	if cfgEnv.Address != "" {
		inst.ServerAddress = cfgEnv.Address
	} else {
		inst.ServerAddress = cfgFromFlags.Address
	}
	inst.Sugar.Infof("Server Strated with variables: address=%v", inst.ServerAddress)

	if cfgFromJsonFile.StoreInterval >= 0 {
		inst.StoreInterval = cfgFromJsonFile.StoreInterval
	}
	if cfgEnv.StoreInterval >= 0 {
		inst.StoreInterval = cfgEnv.StoreInterval
	} else {
		inst.StoreInterval = cfgFromFlags.StoreInterval
	}
	inst.Sugar.Infof("Server Strated with variables: StoreInterval=%v", inst.StoreInterval)

	if len(cfgFromJsonFile.StoreFile) > 0 {
		inst.StoreFile = cfgFromJsonFile.StoreFile
	}
	if len(cfgEnv.StoreFile) > 0 {
		inst.StoreFile = cfgEnv.StoreFile
	} else {
		inst.StoreFile = cfgFromFlags.StoreFile
	}
	inst.Sugar.Infof("Server Strated with variables: StoreFile=%v", inst.StoreFile)

	if cfgFromJsonFile.Key != "" {
		inst.Key = cfgFromJsonFile.Key
	}
	if cfgEnv.Key != "" {
		inst.Key = cfgEnv.Key
	} else {
		inst.Key = cfgFromFlags.Key
	}
	inst.Sugar.Infof("Server Strated with variables: Key=%v", inst.Key)

	if cfgFromJsonFile.Restore {
		inst.Restore = true
	}
	if cfgEnv.Restore {
		inst.Restore = true
	} else if cfgFromFlags.Restore {
		inst.Restore = true
	} else {
		inst.Restore = false
	}
	inst.Sugar.Infof("Server Strated with variables: Restore=%v", inst.Restore)

	if cfgFromJsonFile.DatabaseDSN != "" {
		inst.DatabaseDSN = cfgFromJsonFile.DatabaseDSN
	}
	if cfgEnv.DatabaseDSN != "" {
		inst.DatabaseDSN = cfgEnv.DatabaseDSN
	} else if cfgFromFlags.DatabaseDSN != "" {
		inst.DatabaseDSN = cfgFromFlags.DatabaseDSN
	}
	inst.Sugar.Infof("Server Strated with variables: DatabaseDSN=%v", inst.DatabaseDSN)

	if cfgFromJsonFile.LogHTTP {
		inst.LogHTTP = true
	}
	if cfgEnv.LogHTTP {
		inst.LogHTTP = true
	} else if cfgFromFlags.LogHTTP {
		inst.LogHTTP = true
	} else {
		inst.LogHTTP = false
	}
	inst.Sugar.Infof("Server Strated with variables: LogHTTP=%v", inst.LogHTTP)

	if cfgFromJsonFile.PrivateKeyFileName != "" {
		inst.PrivateKeyFileName = cfgFromJsonFile.PrivateKeyFileName
	}
	if cfgEnv.PrivateKeyFileName != "" {
		inst.PrivateKeyFileName = cfgEnv.PrivateKeyFileName
	} else if cfgFromFlags.PrivateKeyFileName != "" {
		inst.PrivateKeyFileName = cfgFromFlags.PrivateKeyFileName
	}
	inst.Sugar.Infof("Server Strated with variables: PrivateKeyFileName=%v", inst.PrivateKeyFileName)
}
