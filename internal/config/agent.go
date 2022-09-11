// Package config fill config vars from environment and command flags
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

type AgentConfig struct {
	Address           string        `env:"ADDRESS" json:"address"`
	ReportInterval    time.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval      time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
	Key               string        `env:"KEY" json:"key"`
	PublicKeyFileName string        `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile        string        `env:"CONFIG"`
}

func ConfigAgentInit() {
	var cfgEnv AgentConfig
	var cfgFromFlags AgentConfig
	var cfgFromJsonFile AgentConfig

	// загрузим переменные среды
	err := env.Parse(&cfgEnv)
	if err != nil {
		log.Fatal(err)
	}

	inst.Sugar.Infow("Agent started")
	inst.Sugar.Infof("Config environment:%+v", cfgEnv)

	// load flags
	flag.StringVar(&cfgFromFlags.Address, "a", inst.ServerAddress, "address to bind on")
	flag.StringVar(&cfgFromFlags.Key, "k", "", "key for hash calculate")
	flag.DurationVar(&cfgFromFlags.ReportInterval, "r", inst.ReportInterval, "report interval")
	flag.DurationVar(&cfgFromFlags.PollInterval, "p", inst.PollInterval, "poll interval")
	flag.StringVar(&cfgFromFlags.PublicKeyFileName, "crypto-key", inst.PublicKeyFileName, "certificate with public key file name")
	flag.StringVar(&cfgFromFlags.ConfigFile, "c", "", "config file name")

	flag.BoolVar(&inst.BatchSend, "B", true, "batch send data")
	flag.Parse()
	inst.Sugar.Infof("Agent Config flags:%+v", cfgFromFlags)

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
		inst.ServerAddress = cfgFromFlags.Address
	}
	if cfgEnv.Address != "" {
		inst.ServerAddress = cfgEnv.Address
	} else {
		inst.ServerAddress = cfgFromFlags.Address
	}

	if cfgFromJsonFile.Key != "" {
		inst.Key = cfgFromJsonFile.Key
	}
	if cfgEnv.Key != "" {
		inst.Key = cfgEnv.Key
	} else {
		inst.Key = cfgFromFlags.Key
	}

	if cfgFromJsonFile.PollInterval > 0 {
		inst.PollInterval = cfgFromJsonFile.PollInterval
	}
	if cfgEnv.PollInterval > 0 {
		inst.PollInterval = cfgEnv.PollInterval
	} else {
		inst.PollInterval = cfgFromFlags.PollInterval
	}

	if cfgFromJsonFile.ReportInterval > 0 {
		inst.ReportInterval = cfgFromJsonFile.ReportInterval
	}
	if cfgEnv.ReportInterval > 0 {
		inst.ReportInterval = cfgEnv.ReportInterval
	} else {
		inst.ReportInterval = cfgFromFlags.ReportInterval
	}

	if cfgFromJsonFile.PublicKeyFileName != "" {
		inst.PublicKeyFileName = cfgFromJsonFile.PublicKeyFileName
	}
	if cfgEnv.PublicKeyFileName != "" {
		inst.PublicKeyFileName = cfgEnv.PublicKeyFileName
	} else {
		inst.PublicKeyFileName = cfgFromFlags.PublicKeyFileName
	}

	inst.Sugar.Infof("Agent Strated with variables: address=%v, poll interval=%v, report interval=%v, key=%v, batch send=%v, PublicKeyFileName=%s, ConfigFileName=%s",
		inst.ServerAddress, inst.PollInterval, inst.ReportInterval, inst.Key, inst.BatchSend, inst.PublicKeyFileName, ConfigFileName)
}
