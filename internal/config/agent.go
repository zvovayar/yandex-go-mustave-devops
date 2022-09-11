// Package config fill config vars from environment and command flags
package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type AgentConfig struct {
	Address           string        `env:"ADDRESS"`
	ReportInterval    time.Duration `env:"REPORT_INTERVAL"`
	PollInterval      time.Duration `env:"POLL_INTERVAL"`
	Key               string        `env:"KEY"`
	PublicKeyFileName string        `env:"CRYPTO_KEY"`
}

func ConfigAgentInit() {
	var cfgEnv AgentConfig
	var cfgFromFlags AgentConfig
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

	flag.BoolVar(&inst.BatchSend, "B", true, "batch send data")
	flag.Parse()
	inst.Sugar.Infof("Agent Config flags:%+v", cfgFromFlags)

	// assign work parameters
	if cfgEnv.Address != "" {
		inst.ServerAddress = cfgEnv.Address
	} else {
		inst.ServerAddress = cfgFromFlags.Address
	}
	if cfgEnv.Key != "" {
		inst.Key = cfgEnv.Key
	} else {
		inst.Key = cfgFromFlags.Key
	}
	if cfgEnv.PollInterval > 0 {
		inst.PollInterval = cfgEnv.PollInterval
	} else {
		inst.PollInterval = cfgFromFlags.PollInterval
	}
	if cfgEnv.ReportInterval > 0 {
		inst.ReportInterval = cfgEnv.ReportInterval
	} else {
		inst.ReportInterval = cfgFromFlags.ReportInterval
	}
	if cfgEnv.PublicKeyFileName != "" {
		inst.PublicKeyFileName = cfgEnv.PublicKeyFileName
	} else {
		inst.PublicKeyFileName = cfgFromFlags.PublicKeyFileName
	}

	inst.Sugar.Infof("Agent Strated with variables: address=%v, poll interval=%v, report interval=%v, key=%v, batch send=%v, PublicKeyFileName=%s",
		inst.ServerAddress, inst.PollInterval, inst.ReportInterval, inst.Key, inst.BatchSend, inst.PublicKeyFileName)
}
