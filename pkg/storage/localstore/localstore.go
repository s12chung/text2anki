// Package localstore stores file locally, tries to mimic s3
package localstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

// API has the API to do file storage
type API struct {
	origin    string
	basePath  string
	encryptor Encryptor
}

// NewAPI returns a new API
func NewAPI(origin, basePath string, encryptor Encryptor) API {
	if !strings.HasSuffix(origin, "/") {
		origin += "/"
	}
	return API{origin: origin, basePath: basePath, encryptor: encryptor}
}

// CipherQueryParam is the query parameter that contains the ciphertext for signing
const CipherQueryParam = "ciphertext"

// Sign returns a storage.PresignedHTTPRequest
func (a API) Sign(key string) (storage.PresignedHTTPRequest, error) {
	ciphertext, err := a.encryptor.Encrypt(key)
	if err != nil {
		return storage.PresignedHTTPRequest{}, err
	}
	return storage.PresignedHTTPRequest{
		URL:          a.origin + key + "?" + url.Values{CipherQueryParam: []string{ciphertext}}.Encode(),
		Method:       "PUT",
		SignedHeader: nil,
	}, nil
}

// Validate validates whether the given key and values match for signing
func (a API) Validate(key string, values url.Values) error {
	ciphertext := values.Get(CipherQueryParam)
	cipherKey, err := a.encryptor.Decrypt(ciphertext)
	if err != nil {
		return err
	}
	if cipherKey != key {
		return fmt.Errorf("ciphertext (%v) does not match key (%v)", ciphertext, key)
	}
	return nil
}

// Store stores the file at key, checking if it was signed from the values
func (a API) Store(key string, file io.Reader) error {
	p := path.Join(a.basePath, key)
	if err := os.MkdirAll(filepath.Dir(p), ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}
	outFile, err := os.Create(p) //nolint:gosec //it's the purpose of this package
	if err != nil {
		return err
	}
	if _, err = io.Copy(outFile, file); err != nil {
		return err
	}

	return outFile.Close()
}

// Encryptor defines the interface to encrypt and sign keys, ciphers are base64.URLEncoded
type Encryptor interface {
	Encrypt(message string) (string, error)
	Decrypt(cipher string) (string, error)
}

// AESEncryptor is an AES encryptor
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor returns a new AESEncryptor from the key
func NewAESEncryptor(key []byte) AESEncryptor {
	return AESEncryptor{key: key}
}

// NewAESEncryptorFromFile returns a new AESEncryptor from the keyFile
// generate file via `openssl rand -hex 32`
func NewAESEncryptorFromFile(keyFile string) (AESEncryptor, error) {
	file, err := os.ReadFile(keyFile) //nolint:gosec // want to read the keyfile, read internally
	if err != nil {
		return AESEncryptor{}, err
	}
	key, err := hex.DecodeString(string(file))
	if err != nil {
		return AESEncryptor{}, err
	}
	return AESEncryptor{key: key}, nil
}

// Encrypt encrypts the message
func (a AESEncryptor) Encrypt(message string) (string, error) {
	byteMessage := []byte(message)

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(byteMessage))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:], byteMessage)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the ciphertext
func (a AESEncryptor) Decrypt(ciphertext string) (string, error) {
	cipherDecoded, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	if len(cipherDecoded) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	iv := cipherDecoded[:aes.BlockSize]
	cipherDecoded = cipherDecoded[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherDecoded, cipherDecoded)
	return string(cipherDecoded), nil
}
