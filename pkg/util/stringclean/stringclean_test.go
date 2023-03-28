package stringclean

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpeaker(t *testing.T) {
	tcs := []struct {
		name     string
		s        string
		expected string
	}{
		{name: "none", s: "You should definitely get one.", expected: "You should definitely get one."},
		{name: "simple", s: "Kyeong-Eun: You should definitely get one.", expected: "You should definitely get one."},
		{name: "time", s: "At 3:30, you should definitely get one.", expected: "At 3:30, you should definitely get one."},
		{name: "korean", s: "경은: 나중에 꼭 한번 키워 보세요.", expected: "나중에 꼭 한번 키워 보세요."},
		{name: "long", s: "Cheong Kyeong-Eunnie-Ya: You should definitely get one.", expected: "You should definitely get one."},
		{name: "too long",
			s:        "Cheong Cheong Kyeong-Eunnie-Ya: You should definitely get one.",
			expected: "Cheong Cheong Kyeong-Eunnie-Ya: You should definitely get one."},
		{name: "broken", s: ":      You should definitely get one.    ", expected: "You should definitely get one."},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, Speaker(tc.s))
		})
	}
}
