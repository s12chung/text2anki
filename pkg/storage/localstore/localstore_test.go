package localstore

import (
	"bytes"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var apiOrigin = "http://localhost:3000"

func testAPI(t *testing.T) API {
	return NewAPI(apiOrigin, path.Join(os.TempDir(), test.GenerateName("filestore-api")), testEncryptor(t))
}

func testEncryptor(t *testing.T) AESEncryptor {
	require := require.New(t)
	encryptor, err := NewAESEncryptorFromFile(fixture.JoinTestData("test.key"))
	require.NoError(err)
	return encryptor
}

const testKeyPrefix = "some_table_name/my_columns_me_now/123e4567-e89b-12d3-a456-426614174000"
const testKeyFile = "0.txt"
const testKey = testKeyPrefix + "/" + testKeyFile

func TestAPI_SignPut(t *testing.T) {
	require := require.New(t)

	api := testAPI(t)
	req, err := api.SignPut(testKey)
	require.NoError(err)

	require.Equal("PUT", req.Method)
	require.Empty(req.SignedHeader)

	u, err := url.Parse(req.URL)
	require.NoError(err)

	key, err := api.encryptor.Decrypt(u.Query().Get(CipherQueryParam))
	require.NoError(err)
	require.Equal(testKey, key)

	u.RawQuery = ""
	require.Equal(apiOrigin+"/"+testKey, u.String())
}

func TestAPI_SignGet(t *testing.T) {
	prefix := "TestAPI_SignGet/test/me"
	require := require.New(t)

	api := testAPI(t)
	u, err := api.SignGet(prefix + "/" + testKeyFile)
	require.Equal(errSignGetNotFound, err)
	require.Empty(u)

	require.NoError(api.Store(prefix+"/"+testKeyFile, bytes.NewReader([]byte("test_me"))))
	u, err = api.SignGet(prefix + "/" + testKeyFile)
	require.NoError(err)
	require.Equal(api.keyURL(prefix+"/"+testKeyFile), u)
}

func TestAPI_ListKeys(t *testing.T) {
	prefix := "TestAPI_ListKeys/test/me"
	require := require.New(t)

	api := testAPI(t)
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

func TestAPI_Validate(t *testing.T) {
	require := require.New(t)
	api := testAPI(t)
	ciphertext, err := api.encryptor.Encrypt(testKey)
	require.NoError(err)
	require.NoError(api.Validate(testKey, url.Values{CipherQueryParam: []string{ciphertext}}))
	require.Error(api.Validate(testKey, url.Values{}))
	require.Error(api.Validate(testKey, url.Values{CipherQueryParam: []string{"bad_cipher"}}))
}

func TestAPI_Store(t *testing.T) {
	testStore(t)
	testStore(t)
}

func testStore(t *testing.T) {
	require := require.New(t)

	api := testAPI(t)
	fileData := []byte("Store")
	require.NoError(api.Store(testKey, bytes.NewReader(fileData)))

	fileBytes, err := os.ReadFile(path.Join(api.keyBasePath, testKey))
	require.NoError(err)
	require.Equal(fileData, fileBytes)
}

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	require := require.New(t)

	cipher, err := testEncryptor(t).Encrypt(testKey)
	require.NoError(err)
	message, err := testEncryptor(t).Decrypt(cipher)
	require.NoError(err)
	require.Equal(testKey, message)
}
