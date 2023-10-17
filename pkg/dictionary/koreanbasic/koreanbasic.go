// Package koreanbasic contains functions for the Korean Basic Dictionary API
package koreanbasic

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
)

// DictionarySource is the name of the dictionary
const DictionarySource = "Korean Basic Dictionary"

// GetAPIKeyFromEnv gets the API key from the default ENV var
func GetAPIKeyFromEnv() string {
	return os.Getenv("KOREAN_BASIC_API_KEY")
}

// KoreanBasic is a Korean Basic dictionary API wrapper
type KoreanBasic struct {
	apiKey string

	client *http.Client
}

// New returns a KoreanBasic dictionary
func New(apiKey string) *KoreanBasic {
	transport := &http.Transport{
		//nolint:gosec // needed
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &KoreanBasic{apiKey: apiKey, client: &http.Client{Transport: transport}}
}

// Search returns the search results of the query
func (k *KoreanBasic) Search(ctx context.Context, q string, pos lang.PartOfSpeech) ([]dictionary.Term, error) {
	bytes, err := k.getSearch(ctx, q, pos)
	if err != nil {
		return nil, err
	}
	channel, err := unmarshallSearch(bytes)
	if err != nil {
		return nil, err
	}
	return itemsToTerms(channel.Items)
}

func (k *KoreanBasic) getSearch(ctx context.Context, q string, pos lang.PartOfSpeech) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL(k.apiKey, q, pos), nil)
	if err != nil {
		return nil, err
	}
	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck // failing is ok
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response from API is not OK (200), got: %v (%v)", resp.Status, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// SetClient sets the client for API requests
func (k *KoreanBasic) SetClient(c *http.Client) {
	k.client = c
}

const apiURLStart = "https://krdict.korean.go.kr/api/search?sort=popular&translated=y&trans_lang=1"
const apiURLAdvanced = "&advanced=y&method=include"
const apiURLKeyQuery = "&key=%s&q=%s"
const apiURLPos = "&pos=%v"

func apiURL(apiKey, q string, pos lang.PartOfSpeech) string {
	pos = mergePosMap[pos]

	urlString := apiURLStart
	posString, exists := partOfSpeechReverseMap[pos]
	if exists {
		urlString += apiURLAdvanced
	}
	urlString += fmt.Sprintf(apiURLKeyQuery, apiKey, url.QueryEscape(q))
	if exists {
		urlString += fmt.Sprintf(apiURLPos, partOfSpeechToAPIInt[posString])
	}
	return urlString
}

func unmarshallSearch(data []byte) (*channel, error) {
	ch := &channel{}
	if err := xml.Unmarshal(data, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

var wordGradeToCommonLevel = map[string]lang.CommonLevel{
	"":   lang.CommonLevelUnique,
	"고급": lang.CommonLevelRare,
	"중급": lang.CommonLevelMedium,
	"초급": lang.CommonLevelCommon,
}

var partOfSpeechMap = map[string]lang.PartOfSpeech{
	"명사":  lang.PartOfSpeechNoun,
	"대명사": lang.PartOfSpeechPronoun,
	"수사":  lang.PartOfSpeechNumeral,
	"조사":  lang.PartOfSpeechPostposition,

	"동사":  lang.PartOfSpeechVerb,
	"형용사": lang.PartOfSpeechAdjective,
	"관형사": lang.PartOfSpeechDeterminer,

	"부사":  lang.PartOfSpeechAdverb,
	"감탄사": lang.PartOfSpeechInterjection,

	"접사": lang.PartOfSpeechAffix,

	"의존 명사": lang.PartOfSpeechDependentNoun,

	"보조 동사":  lang.PartOfSpeechAuxiliaryVerb,
	"보조 형용사": lang.PartOfSpeechAuxiliaryAdjective,

	"어미": lang.PartOfSpeechEnding,

	"품사 없음": lang.PartOfSpeechUnknown,
	"":      lang.PartOfSpeechUnknown,
}

var mergePosMap = map[lang.PartOfSpeech]lang.PartOfSpeech{
	lang.PartOfSpeechNoun:         lang.PartOfSpeechNoun,
	lang.PartOfSpeechPronoun:      lang.PartOfSpeechPronoun,
	lang.PartOfSpeechNumeral:      lang.PartOfSpeechNumeral,
	lang.PartOfSpeechAlphabet:     lang.PartOfSpeechNoun, // Dictionary Encoding
	lang.PartOfSpeechPostposition: lang.PartOfSpeechPostposition,

	lang.PartOfSpeechVerb:       lang.PartOfSpeechVerb,
	lang.PartOfSpeechAdjective:  lang.PartOfSpeechAdjective,
	lang.PartOfSpeechDeterminer: lang.PartOfSpeechDeterminer,

	lang.PartOfSpeechAdverb:       lang.PartOfSpeechAdverb,
	lang.PartOfSpeechInterjection: lang.PartOfSpeechInterjection,

	lang.PartOfSpeechAffix:  lang.PartOfSpeechAffix,
	lang.PartOfSpeechPrefix: lang.PartOfSpeechAffix, // Make them the same
	lang.PartOfSpeechInfix:  lang.PartOfSpeechAffix, // Make them the same
	lang.PartOfSpeechSuffix: lang.PartOfSpeechAffix, // Make them the same

	lang.PartOfSpeechRoot: lang.PartOfSpeechEnding, // Convert, but untested

	lang.PartOfSpeechDependentNoun: lang.PartOfSpeechDependentNoun,

	lang.PartOfSpeechAuxiliaryPredicate: lang.PartOfSpeechUnknown, // Convert
	lang.PartOfSpeechAuxiliaryVerb:      lang.PartOfSpeechAuxiliaryVerb,
	lang.PartOfSpeechAuxiliaryAdjective: lang.PartOfSpeechAuxiliaryAdjective,

	lang.PartOfSpeechEnding:      lang.PartOfSpeechEnding,
	lang.PartOfSpeechCopula:      lang.PartOfSpeechPostposition, // Convert
	lang.PartOfSpeechPunctuation: lang.PartOfSpeechEmpty,        // Skip

	lang.PartOfSpeechOtherLanguage: lang.PartOfSpeechEmpty, // Skip
	lang.PartOfSpeechOther:         lang.PartOfSpeechEmpty, // Skip
	lang.PartOfSpeechUnknown:       lang.PartOfSpeechUnknown,
	lang.PartOfSpeechEmpty:         lang.PartOfSpeechEmpty,
}

var partOfSpeechReverseMap map[lang.PartOfSpeech]string

func init() {
	partOfSpeechReverseMap = map[lang.PartOfSpeech]string{}
	for k, v := range partOfSpeechMap {
		if v != lang.PartOfSpeechUnknown {
			partOfSpeechReverseMap[v] = k
		}
	}
}

var partOfSpeechToAPIInt = map[string]uint{
	"명사":     1,
	"대명사":    2,
	"수사":     3,
	"조사":     4,
	"동사":     5,
	"형용사":    6,
	"관형사":    7,
	"부사":     8,
	"감탄사":    9,
	"접사":     10,
	"의존 명사":  11,
	"보조 동사":  12,
	"보조 형용사": 13,
	"어미":     14,
}

func itemsToTerms(items []item) ([]dictionary.Term, error) {
	terms := make([]dictionary.Term, 0, len(items))
	for _, item := range items {
		term, err := item.term()
		if err != nil {
			return nil, err
		}
		if term.Text == "" {
			continue
		}
		terms = append(terms, term)
	}
	return terms, nil
}

type channel struct {
	Items  []item `xml:"item"`
	Total  uint   `xml:"total"`
	Start  uint   `xml:"start"`
	Number uint   `xml:"num"`
}

type item struct {
	TargetCode   uint    `xml:"target_code"`
	Word         string  `xml:"word"`
	WordGrade    string  `xml:"word_grade"`
	PartOfSpeech string  `xml:"pos"`
	Senses       []sense `xml:"sense"`
}

func (i *item) term() (dictionary.Term, error) {
	if _, exists := partOfSpeechMap[i.PartOfSpeech]; !exists {
		return dictionary.Term{}, fmt.Errorf("part of speech not found: %v, %v", i.Word, i.PartOfSpeech)
	}
	if len(i.Senses) == 0 {
		return dictionary.Term{}, nil
	}

	term := dictionary.Term{
		ID:   int64(i.TargetCode),
		Text: strings.TrimSpace(i.Word),
		// Not supported by API, try for karaoke, "가라오케", variants: "가라오께, 가라오게, 까라오께, 카라오케, 카라오께"
		Variants:         []string{},
		CommonLevel:      wordGradeToCommonLevel[i.WordGrade],
		PartOfSpeech:     partOfSpeechMap[i.PartOfSpeech],
		DictionarySource: DictionarySource,
	}
	term.Translations = i.translations()
	return term, nil
}

func (i *item) translations() []dictionary.Translation {
	translations := make([]dictionary.Translation, len(i.Senses))
	for j, sense := range i.Senses {
		translations[j] = sense.translation()
	}
	return translations
}

type sense struct {
	Translation string `xml:"translation>trans_word"`
	Explanation string `xml:"translation>trans_dfn"`
}

func (s *sense) translation() dictionary.Translation {
	return dictionary.Translation{
		Text:        strings.TrimSpace(s.Translation),
		Explanation: strings.TrimSpace(s.Explanation),
	}
}
