package storage

import (
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

func (t testAPI) SignPut(key string) (PresignedHTTPRequest, error) {
	key = uuidRegexp.ReplaceAllString(key, testUUID)
	return PresignedHTTPRequest{
		URL:          path.Join("http://localhost:3000", key) + "?cipher=blah",
		Method:       "PUT",
		SignedHeader: nil,
	}, nil
}

func TestSigner_SignPut(t *testing.T) {
	require := require.New(t)
	testName := "TestSigner_SignPut"

	reqs, id, err := NewSigner(testAPI{}).SignPut("sources", "parts", []string{".jpg", ".png", ".jpeg"})
	require.NoError(err)
	require.NotEqual(testUUID, id)
	for _, req := range reqs {
		require.Contains(req.URL, testUUID)
	}
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, reqs))
}
