package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	oldTime := timeNow
	formattedTime := int64(1605139200) // This is 2020-11-12 00:00:00 +0000 UTC
	timeNow = func() time.Time {
		return time.Unix(formattedTime, 0)
	}
	exit := m.Run()
	timeNow = oldTime
	os.Exit(exit)
}

func TestGenerateFilename(t *testing.T) {
	expected := fmt.Sprintf("text2anki-waka-%v.blah", timeNow().Format(time.StampNano))

	tcs := []struct {
		name     string
		ext      string
		expected string
	}{
		{name: "waka", ext: ".blah", expected: expected},
		{name: "waka", ext: "blah", expected: expected},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, GenerateFilename(tc.name, tc.ext))
		})
	}
}

func TestGenerateName(t *testing.T) {
	require := require.New(t)
	expected := fmt.Sprintf("text2anki-waka-%v", timeNow().Format(time.StampNano))
	require.Equal(expected, GenerateName("waka"))
}
