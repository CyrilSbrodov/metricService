package crypto

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
)

type Cryptoer interface {
	AddCryptoKey(filenamePublicKey, filenamePrivateKey, filenameCert string) error
	CreateNewCryptoFile(PEM bytes.Buffer, filename string) error
	DecryptedData(b []byte, privateKey *rsa.PrivateKey) ([]byte, error)
	EncryptedData(b []byte, publicKey *rsa.PublicKey) ([]byte, error)
	LoadPrivatePEMKey(filename string) (*rsa.PrivateKey, error)
	LoadPublicPEMKey(filename string) (*rsa.PublicKey, error)
}

type Crypto struct {
	logger loggers.Logger
}

func NewCrypto() *Crypto {
	logger := loggers.NewLogger()
	return &Crypto{
		logger: *logger,
	}
}

func (c *Crypto) AddCryptoKey(filenamePublicKey, filenamePrivateKey, filenameCert string) error {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"metricService"},
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
		c.logger.LogErr(err, "")
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		c.logger.LogErr(err, "")
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

	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	if err = c.CreateNewCryptoFile(certPEM, filenameCert); err != nil {
		c.logger.LogErr(err, "filed to create new file")
		return err
	}
	if err = c.CreateNewCryptoFile(privateKeyPEM, filenamePrivateKey); err != nil {
		c.logger.LogErr(err, "filed to create new file")
		return err
	}
	if err = c.CreateNewCryptoFile(publicKeyPEM, filenamePublicKey); err != nil {
		c.logger.LogErr(err, "filed to create new file")
	}
	return nil
}

func (c *Crypto) CreateNewCryptoFile(PEM bytes.Buffer, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	//"../../internal/crypto/"+
	if err != nil {
		c.logger.LogErr(err, "failed to open/create file")
		return err
	}
	writer := bufio.NewWriter(file)
	if _, err := writer.Write(PEM.Bytes()); err != nil {
		c.logger.LogErr(err, "failed to write buffer")
		return err
	}
	writer.Flush()
	defer file.Close()
	return nil
}

//func (c *Crypto) DecryptedData(b []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
//	label := []byte("OAEP Encrypted")
//	rng := rand.Reader
//	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, privateKey, b, label)
//	if err != nil {
//		c.logger.LogErr(err, "error from decryption")
//		return nil, err
//	}
//	return plaintext, nil
//}

func (c *Crypto) DecryptedData(msg []byte, private *rsa.PrivateKey) ([]byte, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(sha256.New(), rng, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func (c *Crypto) EncryptedData(b []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	label := []byte("OAEP Encrypted")

	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, b, label)
	if err != nil {
		c.logger.LogErr(err, "error from encryption")
	}
	return ciphertext, nil
}

func (c *Crypto) LoadPrivatePEMKey(filename string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(filename)
	//"../../internal/crypto/" +
	if err != nil {
		c.logger.LogErr(err, "filed to open file")
		return nil, err
	}
	defer privateKeyFile.Close()

	pemFileInfo, err := privateKeyFile.Stat()
	if err != nil {
		c.logger.LogErr(err, "filed to read stat from file")
	}
	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)

	if err != nil {
		c.logger.LogErr(err, "filed to read bytes from file")
		return nil, err
	}

	data, _ := pem.Decode(pemBytes)

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		c.logger.LogErr(err, "filed to import key")
		return nil, err
	}
	return privateKeyImported, nil
}

func (c *Crypto) LoadPublicPEMKey(filename string) (*rsa.PublicKey, error) {
	publicKeyFile, err := os.Open(filename)
	//"../../internal/crypto/" +
	if err != nil {
		c.logger.LogErr(err, "filed to open file")
		return nil, err
	}
	defer publicKeyFile.Close()

	pemFileInfo, err := publicKeyFile.Stat()
	if err != nil {
		c.logger.LogErr(err, "filed to read stat from file")
	}
	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pemBytes)

	if err != nil {
		c.logger.LogErr(err, "filed to read bytes from file")
		return nil, err
	}

	data, _ := pem.Decode(pemBytes)

	public, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		c.logger.LogErr(err, "filed to import key")
		return nil, err
	}
	return public, nil
}

func LoadPublicPEMKey(filename string, logger *loggers.Logger) (*rsa.PublicKey, error) {
	publicKeyFile, err := os.Open("../../cmd/server/" + filename)
	//"../../internal/crypto/" +
	if err != nil {
		logger.LogErr(err, "filed to open file")
		return nil, err
	}
	defer publicKeyFile.Close()

	pemFileInfo, err := publicKeyFile.Stat()
	if err != nil {
		logger.LogErr(err, "filed to read stat from file")
	}
	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pemBytes)

	if err != nil {
		logger.LogErr(err, "filed to read bytes from file")
		return nil, err
	}

	data, _ := pem.Decode(pemBytes)

	public, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		logger.LogErr(err, "filed to import key")
		return nil, err
	}
	return public, nil
}

//func EncryptedData(b []byte, publicKey *rsa.PublicKey, logger *loggers.Logger) ([]byte, error) {
//	label := []byte("OAEP Encrypted")
//
//	rng := rand.Reader
//	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, b, label)
//	if err != nil {
//		logger.LogErr(err, "error from encryption")
//	}
//	return ciphertext, nil
//}

func EncryptedData(msg []byte, public *rsa.PublicKey, logger *loggers.Logger) ([]byte, error) {
	label := []byte("OAEP Encrypted")

	rng := rand.Reader

	msgLen := len(msg)
	step := public.Size() - 2*sha256.New().Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(sha256.New(), rng, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}
