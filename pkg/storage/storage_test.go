package storage

import (
	"net/http"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var testUUID = "123e4567-e89b-12d3-a456-426614174000"
var uuidRegexp = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)

type testAPI struct {
}

func keyURL(key string) string {
	return "http://localhost:3000/" + key
}

func (t testAPI) SignPut(key string) (PreSignedHTTPRequest, error) {
	key = uuidRegexp.ReplaceAllString(key, testUUID)
	return PreSignedHTTPRequest{
		URL:          keyURL(key) + "?cipher=blah",
		Method:       "PUT",
		SignedHeader: http.Header{},
	}, nil
}

func (t testAPI) SignGet(key string) (string, error) {
	return keyURL(key), nil
}

func (t testAPI) ListKeys(prefix string) ([]string, error) {
	return []string{path.Join(prefix, "a.txt"), path.Join(prefix, "b.txt")}, nil
}

func TestSigner_SignGetByID(t *testing.T) {
	require := require.New(t)
	urls, err := NewSigner(testAPI{}).SignGetByID("sources", "parts", testUUID)
	require.NoError(err)

	prefix := "http://localhost:3000/sources/parts/123e4567-e89b-12d3-a456-426614174000/"
	require.Equal([]string{prefix + "a.txt", prefix + "b.txt"}, urls)
}

func testSignPutBuilder() signPutBuilder {
	return newSignPutBuilder("sources", "parts", testUUID, testAPI{})
}

func TestSignPutBuilder_ID(t *testing.T) {
	require := require.New(t)
	require.Equal(testUUID, testSignPutBuilder().ID())
}

func TestSignPutBuilder_Index_Field_Sign(t *testing.T) {
	require := require.New(t)
	testName := "TestSignPutBuilder_Index_Field_Sign"

	req, err := testSignPutBuilder().Index(1).Field("field1").Field("field2").Index(2).Sign(".txt")
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, req))
}
