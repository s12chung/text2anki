package localstore

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testKeyPrefix = "some_table_name/my_columns_me_now/123e4567-e89b-12d3-a456-426614174000"
const testKeyFile = "0.txt"
const testKey = testKeyPrefix + "/" + testKeyFile

const storagePrefix = "localstore_test"
const apiOrigin = "http://localhost:3000"
const encryptorKey = "4b04fcad4d66da4cb441901cd91e5f508aef00328017ed83f4b0c70fa269a51d" //nolint:gosec // a fake key

func newAPIWithT(t *testing.T) API {
	return NewAPI(apiOrigin, path.Join(os.TempDir(), test.GenerateName(storagePrefix)), newEncryptorT(t))
}

func newEncryptorT(t *testing.T) AESEncryptor {
	require := require.New(t)

	key, err := hex.DecodeString(encryptorKey)
	require.NoError(err)
	encryptor, err := NewAESEncryptor(key)
	require.NoError(err)
	return encryptor
}

func TestAPI_SignPut(t *testing.T) {
	require := require.New(t)

	req, err := newAPIWithT(t).SignPut(testKey)
	require.NoError(err)

	require.Equal("PUT", req.Method)
	require.Empty(req.SignedHeader)

	u, err := url.Parse(req.URL)
	require.NoError(err)

	key, err := newEncryptorT(t).Decrypt(u.Query().Get(CipherQueryParam))
	require.NoError(err)
	require.Equal(testKey, key)

	u.RawQuery = ""
	require.Equal(apiOrigin+"/"+testKey, u.String())
}

func TestAPI_SignGet(t *testing.T) {
	require := require.New(t)

	key := "TestAPI_SignGet/test/me/" + testKeyFile

	api := newAPIWithT(t)
	u, err := api.SignGet(key)
	require.Equal(fmt.Errorf("file does not exist"), err)
	require.Empty(u)

	require.NoError(api.Store(key, bytes.NewReader([]byte("test_me"))))
	u, err = api.SignGet(key)
	require.NoError(err)
	require.Equal(apiOrigin+"/"+key, u)
}

func TestAPI_KeyFromSignGet(t *testing.T) {
	require := require.New(t)

	expectedKey := "TestAPI_SignGet/test/me/" + testKeyFile

	api := newAPIWithT(t)
	require.NoError(api.Store(expectedKey, bytes.NewReader([]byte("test_me"))))
	signGet, err := api.SignGet(expectedKey)
	require.NoError(err)

	key, err := api.KeyFromSignGet(signGet)
	require.NoError(err)
	require.Equal(expectedKey, key)
}

func TestAPI_ValidateAny(t *testing.T) {
	require := require.New(t)
	api := newAPIWithT(t)
	ciphertext, err := newEncryptorT(t).Encrypt(testKey)
	require.NoError(err)
	require.NoError(api.Validate(testKey, url.Values{CipherQueryParam: []string{ciphertext}}))
	require.Error(api.Validate(testKey, url.Values{}))
	require.Error(api.Validate(testKey, url.Values{CipherQueryParam: []string{"bad_cipher"}}))
}

func TestAPI_ListKeys(t *testing.T) {
	require := require.New(t)

	prefix := "TestAPI_ListKeys/test/me"
	api := newAPIWithT(t)
	keys, err := api.ListKeys(prefix)
	require.NoError(err)
	require.Len(keys, 0)

	key1 := path.Join(prefix, testKeyFile)
	key2 := path.Join(prefix, "again_me.txt")
	require.NoError(api.Store(key1, bytes.NewReader([]byte("test_me"))))
	require.NoError(api.Store(key2, bytes.NewReader([]byte("again"))))

	expectedKeys := []string{key1, key2}
	keys, err = api.ListKeys(prefix)
	require.NoError(err)
	require.Equal(expectedKeys, keys)

	keys, err = api.ListKeys(prefix + "/")
	require.NoError(err)
	require.Equal(expectedKeys, keys)
}

func TestAPI_StoreGet(t *testing.T) {
	testStore := func(t *testing.T) {
		require := require.New(t)

		api := newAPIWithT(t)
		fileData := []byte("Store")
		require.NoError(api.Store(testKey, bytes.NewReader(fileData)))

		file, err := api.Get(testKey)
		require.NoError(err)
		fileBytes, err := io.ReadAll(file)
		require.NoError(err)
		require.Equal(fileData, fileBytes)
	}

	testStore(t)
	testStore(t)
}

func TestNewAESEncryptorFromFile(t *testing.T) {
	require := require.New(t)
	testName := "TestNewAESEncryptorFromFile"

	encryptor, err := NewAESEncryptorFromFile(fixture.JoinTestData(testName + ".key"))
	require.NoError(err)

	message := "some message is here not long enough?"
	encrypted, err := encryptor.Encrypt(message)
	require.NoError(err)

	decrypted, err := newEncryptorT(t).Decrypt(encrypted)
	require.NoError(err)
	require.Equal(message, decrypted)
}

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	require := require.New(t)

	encryptor := newEncryptorT(t)

	cipher, err := encryptor.Encrypt(testKey)
	require.NoError(err)
	message, err := encryptor.Decrypt(cipher)
	require.NoError(err)
	require.Equal(testKey, message)
}
