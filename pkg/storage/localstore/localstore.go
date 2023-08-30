// Package localstore stores file locally, tries to mimic s3
package localstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
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
	origin      string
	keyBasePath string
	encryptor   Encryptor
}

// NewAPI returns a new API
func NewAPI(origin, keyBasePath string, encryptor Encryptor) API {
	if !strings.HasSuffix(origin, "/") {
		origin += "/"
	}
	return API{origin: origin, keyBasePath: keyBasePath, encryptor: encryptor}
}

// CipherQueryParam is the query parameter that contains the ciphertext for signing
const CipherQueryParam = "ciphertext"

// SignPut returns a storage.PreSignedHTTPRequest
func (a API) SignPut(key string) (storage.PreSignedHTTPRequest, error) {
	ciphertext, err := a.encryptor.Encrypt(key)
	if err != nil {
		return storage.PreSignedHTTPRequest{}, err
	}
	return storage.PreSignedHTTPRequest{
		URL:          a.keyURL(key) + "?" + url.Values{CipherQueryParam: []string{ciphertext}}.Encode(),
		Method:       "PUT",
		SignedHeader: http.Header{},
	}, nil
}

// SignGet gets the signed URL for the key
func (a API) SignGet(key string) (string, error) {
	p := a.keyPath(key)
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("file does not exist")
	}
	return a.keyURL(key), nil
}

// KeyFromSignGet returns the key given the signGet string
func (a API) KeyFromSignGet(signGet string) (string, error) {
	u, err := url.Parse(signGet)
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	return strings.TrimPrefix(u.String(), a.origin), nil
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

// ListKeys lists the keys for the given path prefix
func (a API) ListKeys(prefix string) ([]string, error) {
	files, err := os.ReadDir(a.keyPath(prefix)) // Replace with your directory
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	keys := make([]string, len(files))
	for i, file := range files {
		keys[i] = path.Join(prefix, file.Name())
	}
	return keys, nil
}

// Store stores the file at key, checking if it was signed from the values
func (a API) Store(key string, file io.Reader) error {
	p := a.keyPath(key)
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

// Get returns the file at key
func (a API) Get(key string) (fs.File, error) {
	return os.Open(a.keyPath(key))
}

// FileHandler returns the http.Handler to serve the files
func (a API) FileHandler() http.Handler {
	return http.FileServer(http.Dir(a.keyBasePath))
}

func (a API) keyPath(key string) string {
	return path.Join(a.keyBasePath, key)
}

func (a API) keyURL(key string) string {
	return a.origin + key
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
