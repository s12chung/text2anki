package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateFilename(t *testing.T) {
	oldTime := timeNow
	unixTime := int64(1605139200) // This is 2020-11-12 00:00:00 +0000 UTC
	timeNow = func() time.Time {
		return time.Unix(unixTime, 0)
	}
	defer func() {
		timeNow = oldTime
	}()
	expected := fmt.Sprintf("text2anki-waka-%v.blah", unixTime)

	tcs := []struct {
		name     string
		ext      string
		expected string
	}{
		{name: "waka", ext: ".blah", expected: expected},
		{name: "waka", ext: "blah", expected: expected},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, GenerateFilename(tc.name, tc.ext))
		})
	}
}
