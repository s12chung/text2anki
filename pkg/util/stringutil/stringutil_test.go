package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	tcs := []struct {
		name     string
		s        string
		expected []string
	}{
		{name: "overall", s: "   test1,   	  test2,test3 ", expected: []string{"test1", "test2", "test3"}},
		{name: "spaces inside", s: " a  test1,   	  test2,test3", expected: []string{"a  test1", "test2", "test3"}},
		{name: "empty", s: "", expected: []string{}},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, SplitClean(tc.s, ","))
		})
	}
}

func TestFirstUnbrokenSubstring(t *testing.T) {
	tcs := []struct {
		name          string
		s             string
		expected      string
		shortExpected string
	}{
		{name: "basic", s: "my name is mario. hihi.", expected: "my name is mario"},
		{name: "dash", s: "Korean Learner - Dictionary", expected: "Korean Learner"},
		{name: "underscore", s: "stringutil_test works!", expected: "stringutil_test works", shortExpected: "stringutil_test"},
		{name: "number", s: "stringutil_test123 works!", expected: "stringutil_test123 works", shortExpected: "stringutil_test123"},
		{name: "symbol", s: "My Corp©  is a good, right?", expected: "My Corp©  is a good"},
		{name: "newline", s: "Here I am\nAgain, here.", expected: "Here I am"},
		{name: "korean", s: "이것은 샘플 파일입니다. tmp/in.txt에 자신의 텍스트를 입력합니다.", expected: "이것은 샘플 파일입니다", shortExpected: "이것은 샘플"},
		{name: "spaces", s: "   my name is      mario    . hihi    .", expected: "my name is      mario"},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, FirstUnbrokenSubstring(tc.s, 100))
			shortExpected := tc.shortExpected
			if shortExpected == "" {
				shortExpected = tc.expected
			}
			require.Equal(shortExpected, FirstUnbrokenSubstring(tc.s, 20))
		})
	}
}

func TestFirstUnbrokenIndex(t *testing.T) {
	tcs := []struct {
		name     string
		s        string
		expected int
	}{
		{name: "basic", s: "my name is mario. hihi.", expected: 16},
		{name: "dash", s: "Korean Learner - Dictionary", expected: 15},
		{name: "em_dash", s: "Korean Learner — Dictionary", expected: 15},
		{name: "underscore", s: "stringutil_test works!", expected: 21},
		{name: "number", s: "stringutil_test123 works!", expected: 24},
		{name: "symbol", s: "My Corp©  is a good, right?", expected: 20},
		{name: "newline", s: "I'm running here again here far away and away\nAgain, here.", expected: 45},
		{name: "korean", s: "이것은 샘플 파일입니다. tmp/in.txt에 자신의 텍스트를 입력합니다.", expected: 32},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, FirstUnbrokenIndex(tc.s))
		})
	}
}
