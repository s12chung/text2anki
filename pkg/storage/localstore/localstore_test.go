package localstore_test

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/storage/localstore/localstoretest"
)

const testKeyPrefix = "some_table_name/my_columns_me_now/123e4567-e89b-12d3-a456-426614174000"
const testKeyFile = "0.txt"
const testKey = testKeyPrefix + "/" + testKeyFile

const testStoragePrefix = "localstore_test"

func TestAPI_SignPut(t *testing.T) {
	require := require.New(t)

	api := localstoretest.NewAPIWithT(t, testStoragePrefix)
	req, err := api.SignPut(testKey)
	require.NoError(err)

	require.Equal("PUT", req.Method)
	require.Empty(req.SignedHeader)

	u, err := url.Parse(req.URL)
	require.NoError(err)

	key, err := localstoretest.NewEncryptorT(t).Decrypt(u.Query().Get(CipherQueryParam))
	require.NoError(err)
	require.Equal(testKey, key)

	u.RawQuery = ""
	require.Equal(localstoretest.APIOrigin+"/"+testKey, u.String())
}

func TestAPI_SignGet(t *testing.T) {
	require := require.New(t)

	key := "TestAPI_SignGet/test/me/" + testKeyFile

	api := localstoretest.NewAPIWithT(t, testStoragePrefix)
	u, err := api.SignGet(key)
	require.Equal(fmt.Errorf("file does not exist"), err)
	require.Empty(u)

	require.NoError(api.Store(key, bytes.NewReader([]byte("test_me"))))
	u, err = api.SignGet(key)
	require.NoError(err)
	require.Equal(localstoretest.APIOrigin+"/"+key, u)
}

func TestAPI_KeyFromSignGet(t *testing.T) {
	require := require.New(t)

	expectedKey := "TestAPI_SignGet/test/me/" + testKeyFile

	api := localstoretest.NewAPIWithT(t, testStoragePrefix)
	require.NoError(api.Store(expectedKey, bytes.NewReader([]byte("test_me"))))
	signGet, err := api.SignGet(expectedKey)
	require.NoError(err)

	key, err := api.KeyFromSignGet(signGet)
	require.NoError(err)
	require.Equal(expectedKey, key)
}

func TestAPI_Validate(t *testing.T) {
	require := require.New(t)
	api := localstoretest.NewAPIWithT(t, testStoragePrefix)
	ciphertext, err := localstoretest.NewEncryptorT(t).Encrypt(testKey)
	require.NoError(err)
	require.NoError(api.Validate(testKey, url.Values{CipherQueryParam: []string{ciphertext}}))
	require.Error(api.Validate(testKey, url.Values{}))
	require.Error(api.Validate(testKey, url.Values{CipherQueryParam: []string{"bad_cipher"}}))
}

func TestAPI_ListKeys(t *testing.T) {
	require := require.New(t)

	prefix := "TestAPI_ListKeys/test/me"
	api := localstoretest.NewAPIWithT(t, testStoragePrefix)
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

		api := localstoretest.NewAPIWithT(t, testStoragePrefix)
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

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	require := require.New(t)

	encryptor := localstoretest.NewEncryptorT(t)

	cipher, err := encryptor.Encrypt(testKey)
	require.NoError(err)
	message, err := encryptor.Decrypt(cipher)
	require.NoError(err)
	require.Equal(testKey, message)
}
