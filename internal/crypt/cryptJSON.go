package crypt

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"os"

	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type Encrypter struct {
	pubKey   *rsa.PublicKey
	fileName string
}

type Decrypter struct {
	privKey  *rsa.PrivateKey
	fileName string
}

// Init load key if nil
func (enc *Encrypter) Init() {
	if enc.pubKey != nil {
		return
	}
	enc.fileName = inst.PublicKeyFileName
	enc.loadPubKey(enc.fileName)
}

// Init load key if nil
func (dec *Decrypter) Init() {
	if dec.privKey != nil {
		return
	}
	dec.fileName = inst.PrivateKeyFileName
	dec.loadPrivKey(dec.fileName)
}

// EncryptBytes encrypt bytes slice
func (enc *Encrypter) EncryptBytes(b []byte) ([]byte, error) {

	if enc.pubKey == nil {
		return nil, errors.New("empty public key ")
	}
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, enc.pubKey, b, nil)
	if err != nil {
		inst.Sugar.Error(err)
		return nil, err
	}

	return ciphertext, nil
}

// DecryptBytes decrypt bytes slice
func (dec *Decrypter) DecryptBytes(b []byte) ([]byte, error) {

	if dec.privKey == nil {
		return nil, errors.New("empty private key ")
	}

	rng := rand.Reader

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, dec.privKey, b, nil)
	if err != nil {
		inst.Sugar.Error(err)
		return nil, err
	}

	return plaintext, nil
}

// loadPrivKey load private key for Decrypter
func (dec *Decrypter) loadPrivKey(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	data := make([]byte, 5000)

	for err != io.EOF {

		_, err = reader.Read(data)
		if err != nil && err != io.EOF {
			inst.Sugar.Error(err)
			return err
		}
	}

	file.Close()

	block, _ := pem.Decode(data)
	if block == nil {
		inst.Sugar.Error("can not find correct data in pem file " + filename)
		return errors.New("can not find correct data in pem file " + filename)
	}

	var privateKeyLoaded *rsa.PrivateKey

	if block.Type == "RSA PRIVATE KEY" {
		privateKeyLoaded, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal("x509.ParsePKCS1PrivateKey(block.Bytes) " + err.Error())
		}
	}

	dec.privKey = privateKeyLoaded

	return nil
}

// loadPubKey load public key for Encrypter
func (enc *Encrypter) loadPubKey(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	data := make([]byte, 5000)

	for err != io.EOF {

		_, err = reader.Read(data)
		if err != nil && err != io.EOF {
			inst.Sugar.Error(err)
			return err
		}
	}

	file.Close()

	block, _ := pem.Decode(data)
	if block == nil {
		inst.Sugar.Error("can not find correct data in pem file " + filename)
		return errors.New("can not find correct data in pem file " + filename)
	}

	var x509certificateLoaded *x509.Certificate

	if block.Type == "CERTIFICATE" {
		x509certificateLoaded, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			inst.Sugar.Fatal("x509.ParseCertificate(block.Bytes) " + err.Error())
			return err
		}
	}

	enc.pubKey = (x509certificateLoaded.PublicKey.(*rsa.PublicKey))
	return nil
}
