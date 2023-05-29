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
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, SplitClean(tc.s, ","))
		})
	}
}
