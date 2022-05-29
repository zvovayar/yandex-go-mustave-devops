package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type HashMetrics interface {
	MakeHashMetrics(key string) string
	ControlHashMetrics(key string) bool
}

type MetricsCrypt struct {
	M   inst.Metrics
	key string
	Msg string
}

func (mc *MetricsCrypt) MakeHashMetrics(key string) string {
	mc.key = key
	var msg string

	if mc.M.MType == "counter" {
		msg = fmt.Sprintf("%s:counter:%d", mc.M.ID, mc.M.Delta)
	} else if mc.M.MType == "gauge" {
		msg = fmt.Sprintf("%s:gauge:%f", mc.M.ID, *mc.M.Value)
	} else {
		return "error"
	}

	mc.Msg = msg

	log.Printf("crypt.MakeHashMetrics msg=%v", msg)

	// keysha256 := sha256.Sum256([]byte(key))
	// keysga256S := keysha256[:]

	h := hmac.New(sha256.New, []byte(mc.key))
	h.Write([]byte(msg))
	sign := h.Sum(nil)

	mc.M.Hash = hex.EncodeToString(sign)
	return mc.M.Hash
}

func (mc *MetricsCrypt) ControlHashMetrics(key string) bool {
	mc.key = key

	var msg string

	if mc.M.MType == "counter" {
		msg = fmt.Sprintf("%s:counter:%d", mc.M.ID, mc.M.Delta)
	} else if mc.M.MType == "gauge" {
		msg = fmt.Sprintf("%s:gauge:%f", mc.M.ID, *mc.M.Value)
	} else {
		return false
	}

	mc.Msg = msg

	log.Printf("crypt.ControlHashMetrics msg=%v", msg)

	// keysha256 := sha256.Sum256([]byte(key))
	// keysga256S := keysha256[:]

	h := hmac.New(sha256.New, []byte(mc.key))
	h.Write([]byte(msg))
	sign := h.Sum(nil)

	data, err := hex.DecodeString(mc.M.Hash)
	if err != nil {
		panic(err)
	}

	if hmac.Equal(sign, data) {
		return true
	}
	// mc.M.Hash = hex.EncodeToString(sign)
	return false
}
