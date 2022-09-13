package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {

	var (
		publicKeyFileName  string
		privateKeyFileName string
	)
	// load flags
	flag.StringVar(&publicKeyFileName, "c", "publickey.pem", "certificate with public key file name")
	flag.StringVar(&privateKeyFileName, "p", "privatekey.pem", "private key file name")
	flag.Parse()

	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// save files with keys
	file, err := os.OpenFile(publicKeyFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	writer.Write(certPEM.Bytes())
	writer.Flush()
	file.Close()

	file, err = os.OpenFile(privateKeyFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		log.Fatal(err)
	}

	writer = bufio.NewWriter(file)

	writer.Write(privateKeyPEM.Bytes())
	writer.Flush()
	file.Close()

	// load certificate with public key from files
	file, err = os.Open(publicKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	data := make([]byte, 5000)

	for err != io.EOF {

		var n int
		n, err = reader.Read(data)
		fmt.Print(" part " + fmt.Sprint(n) + " \n")
		fmt.Print(string(data))

	}

	block, _ := pem.Decode(data)

	var x509certificateLoaded *x509.Certificate

	if block.Type == "CERTIFICATE" {
		x509certificateLoaded, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatal("x509.ParseCertificate(data) " + err.Error())
		}
	}

	file.Close()

	// load private key from files
	file, err = os.Open(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	reader = bufio.NewReader(file)

	data = make([]byte, 5000)

	for err != io.EOF {

		var n int
		n, err = reader.Read(data)
		fmt.Print(" part " + fmt.Sprint(n) + " \n")
		fmt.Print(string(data))

	}

	block, _ = pem.Decode(data)

	var privateKeyLoaded *rsa.PrivateKey

	if block.Type == "RSA PRIVATE KEY" {
		privateKeyLoaded, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal("x509.ParsePKCS1PrivateKey(block.Bytes) " + err.Error())
		}
	}

	// encrypt and decrypt string with loaded keys from files
	fmt.Println("encrypt and decrypt string with loaded keys from files")
	encryptedStr := EncryptString("Secret string", (x509certificateLoaded.PublicKey.(*rsa.PublicKey)))
	DecryptString(encryptedStr, privateKeyLoaded)

	fmt.Println("encrypt and decrypt string with loaded keys from files 2nd")
	encryptedStr = EncryptString("Secret string", (x509certificateLoaded.PublicKey.(*rsa.PublicKey)))
	DecryptString(encryptedStr, privateKeyLoaded)

	// encrypt and decrypt string with original keys
	fmt.Println("encrypt and decrypt string with original keys")
	encryptedStr = EncryptString("Secret string", &privateKey.PublicKey)
	DecryptString(encryptedStr, privateKey)

}

// EncryptString encrypt string with pubKey
func EncryptString(msg string, pubKey *rsa.PublicKey) string {

	secretMessage := []byte(msg)
	label := []byte("orders")

	fmt.Printf("message for encryption: %s\n", msg)

	// crypto/rand.Reader is a good source of entropy for randomizing the
	// encryption function.
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, pubKey, secretMessage, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return "error"
	}

	// Since encryption is a randomized function, ciphertext will be
	// different each time.
	fmt.Printf("EncryptedString: %x\n", ciphertext)
	return hex.EncodeToString(ciphertext)
}

func DecryptString(msg string, prvKey *rsa.PrivateKey) string {

	ciphertext, _ := hex.DecodeString(msg)
	label := []byte("orders")

	fmt.Printf("message for decryption: %s\n", msg)

	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader

	// plaintext, err := prvKey.Decrypt(nil, ciphertext, &rsa.OAEPOptions{
	// 	Hash:  crypto.SHA256,
	// 	Label: label,
	// })
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, prvKey, ciphertext, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		return "error"
	}

	fmt.Printf("Decrypted message: %s\n", string(plaintext))

	return string(plaintext)
	// Remember that encryption only provides confidentiality. The
	// ciphertext should be signed before authenticity is assumed and, even
	// then, consider that messages might be reordered.
}
