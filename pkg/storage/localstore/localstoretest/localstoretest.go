// Package localstoretest provides a testing instances of from the localstore package
package localstoretest

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var keyPath string

func init() {
	_, callerFilePath, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("runtime.Caller not ok for localstoretest")
		os.Exit(-1)
	}
	callerPath := path.Dir(callerFilePath)
	keyPath = path.Join(callerPath, fixture.TestDataDir, "test.key")
}

// APIOrigin returns the Origin used in NewAPIWithT
const APIOrigin = "http://localhost:3000"

// NewAPIWithT returns an API used for testing
func NewAPIWithT(t *testing.T, dirPrefix string) localstore.API {
	return newAPI(dirPrefix, NewEncryptorT(t))
}

// NewEncryptorT returns the encryptor used for testing
func NewEncryptorT(t *testing.T) localstore.AESEncryptor {
	require := require.New(t)
	encryptor, err := NewEncryptor()
	require.NoError(err)
	return encryptor
}

// NewAPI returns an API used for testing
func NewAPI(dirPrefix string) (localstore.API, error) {
	encryptor, err := NewEncryptor()
	if err != nil {
		return localstore.API{}, err
	}
	return newAPI(dirPrefix, encryptor), nil
}

func newAPI(dirPrefix string, encryptor localstore.AESEncryptor) localstore.API {
	return localstore.NewAPI(APIOrigin, path.Join(os.TempDir(), test.GenerateName(dirPrefix)), encryptor)
}

// NewEncryptor returns the encryptor used for testing
func NewEncryptor() (localstore.AESEncryptor, error) {
	return localstore.NewAESEncryptorFromFile(keyPath)
}
