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

	log.Printf("crypt.MakeHashMetrics msg=%v\nmc=%v", msg, mc)

	// h := hmac.New(sha256.New, []byte(mc.key))
	// h.Write([]byte(msg))
	// sign := h.Sum(nil)
	m := []byte(msg)
	k := []byte(key)
	sign := MakeMAC(m, k)

	mc.M.Hash = hex.EncodeToString(sign)
	log.Printf("MakeHashMetrics. sign=%v ", sign)
	log.Printf("crypt.MakeHashMetrics msg=%v\nmc=%v", msg, mc)

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

	log.Printf("crypt.ControlHashMetrics msg=%v\nmc=%v", msg, mc)

	m := []byte(msg)
	mMAC, err := hex.DecodeString(mc.M.Hash)
	if err != nil {
		panic(err)
	}
	k := []byte(key)

	//	mc.M.Hash = "hash" //hex.EncodeToString(sign)
	return ValidMAC(m, mMAC, k)

	// h := hmac.New(sha256.New, []byte(mc.key))
	// h.Write([]byte(msg))
	// sign := h.Sum(nil)

	// data, err := hex.DecodeString(mc.M.Hash)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Printf("ControlHashMetrics. sign=%v data=%v", sign, data)
	// return hmac.Equal(sign, data)

}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	log.Printf("ValidMAC messageMAC=%v\nexpectedMAC=%v", messageMAC, expectedMAC)
	log.Printf("ValidMAC (str)messageMAC=%v\n(str)expectedMAC=%v",
		hex.EncodeToString(messageMAC),
		hex.EncodeToString(expectedMAC))

	return hmac.Equal(messageMAC, expectedMAC)
}

func MakeMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	MAC := mac.Sum(nil)
	log.Printf("MakeMAC expectedMAC=%v", MAC)
	return MAC

}
