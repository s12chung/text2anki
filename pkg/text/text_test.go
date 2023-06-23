package text

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pemistahl/lingua-go"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestLanguagesMatch(t *testing.T) {
	require := require.New(t)
	require.Equal(int(Zulu)+1, len(lingua.AllLanguages()))
	require.Equal(int(Unknown), int(lingua.Unknown))
}

func TestParser_Texts(t *testing.T) {
	testNamePath := "TestParser_Texts/"
	tcs := []struct {
		name string
		err  error
	}{
		{name: "none"},
		{name: "simple_weave"},
		{name: "weave"},
		{name: "weave_skip"},
		{name: "split"},
		{name: "split_1_line"},
		{name: "split_extra_text", err: errExtraTextLine},
		{name: "split_extra_translation", err: errExtraTranslationLine},
		{name: "split_1_line_extra_translation", err: errExtraTranslationLine},
		{name: "weave_extra_translation", err: fmt.Errorf("translation exists for two consecutive non-empty lines: my extra line")},
	}

	parser := NewParser(Korean, English)
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			s := string(fixture.Read(t, testNamePath+tc.name+".txt"))
			split := strings.Split(s, "===")
			if len(split) == 1 {
				split = append(split, "")
			}

			texts, err := parser.Texts(split[0], split[1])
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)

			fixture.CompareReadOrUpdate(t, testNamePath+tc.name+".json", fixture.JSON(t, texts))
		})
	}
}

func TestCleanSpeaker(t *testing.T) {
	texts := CleanSpeaker([]Text{
		{Text: "Mario: It's me", Translation: "경은: 나중에"},
		{Text: "It", Translation: "보세요"},
		{Text: "Bowser: Wit", Translation: "Kyeong-Eunnie-Ya: 나"},
	})
	fixture.CompareReadOrUpdate(t, "TestCleanSpeaker.json", fixture.JSON(t, texts))
}

func TestCleanSpeakerString(t *testing.T) {
	//nolint:dupword
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
			require.Equal(tc.expected, CleanSpeakerString(tc.s))
		})
	}
}
