// Package crypt contain simple crypt functions for validate end create hash
package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

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
		msg = fmt.Sprintf("%s:counter:%d", mc.M.ID, *mc.M.Delta)
	} else if mc.M.MType == "gauge" {
		msg = fmt.Sprintf("%s:gauge:%f", mc.M.ID, *mc.M.Value)
	} else {
		return "error"
	}

	mc.Msg = msg

	// inst.Sugar.Infof("crypt.MakeHashMetrics msg=%v\nmc=%v", msg, mc)

	m := []byte(msg)
	k := []byte(key)
	sign := MakeMAC(m, k)

	mc.M.Hash = hex.EncodeToString(sign)
	// inst.Sugar.Infof("crypt.MakeHashMetrics sign=%v ", sign)
	// inst.Sugar.Infof("crypt.MakeHashMetrics msg=%v\nmc=%v", msg, mc)

	return mc.M.Hash
}

func (mc *MetricsCrypt) ControlHashMetrics(key string) bool {
	mc.key = key

	var msg string

	if mc.M.MType == "counter" {
		msg = fmt.Sprintf("%s:counter:%d", mc.M.ID, *mc.M.Delta)
	} else if mc.M.MType == "gauge" {
		msg = fmt.Sprintf("%s:gauge:%f", mc.M.ID, *mc.M.Value)
	} else {
		return false
	}

	mc.Msg = msg

	// inst.Sugar.Infof("crypt.ControlHashMetrics msg=%v\nmc=%v", msg, mc)

	m := []byte(msg)
	mMAC, err := hex.DecodeString(mc.M.Hash)
	if err != nil {
		panic(err)
	}
	k := []byte(key)

	inst.Sugar.Infof("crypt.ControlHashMetrics m=%v mMAC=%v k=%v",
		string(m),
		hex.EncodeToString(mMAC),
		string(k))
	return ValidMAC(m, mMAC, k)

}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	// inst.Sugar.Infof("ValidMAC messageMAC=%v\nexpectedMAC=%v", messageMAC, expectedMAC)
	// inst.Sugar.Infof("ValidMAC (str)messageMAC=%v\n(str)expectedMAC=%v",
	// hex.EncodeToString(messageMAC),
	// hex.EncodeToString(expectedMAC))

	return hmac.Equal(messageMAC, expectedMAC)
}

// MakeMAC create and return MAC for message with key
func MakeMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	MAC := mac.Sum(nil)
	// inst.Sugar.Infof("MakeMAC expectedMAC=%v", MAC)
	return MAC

}
