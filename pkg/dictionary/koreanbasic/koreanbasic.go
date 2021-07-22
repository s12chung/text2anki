// Package koreanbasic contains functions for the Korean Basic Dictionary API
package koreanbasic

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
)

// DictionarySource is the name of the dictionary
const DictionarySource = "Korean Basic Dictionary"

// KoreanBasic is a Korean Basic dictionary API wrapper
type KoreanBasic struct {
	apiKey string
}

// NewKoreanBasic returns a KoreanBasic dictionary
func NewKoreanBasic(apiKey string) dictionary.Dicionary {
	return &KoreanBasic{apiKey: apiKey}
}

// Search returns the search results of the query
func (k *KoreanBasic) Search(q string) ([]dictionary.Term, error) {
	bytes, err := k.getSearch(q)
	if err != nil {
		return nil, err
	}
	return SearchTerms(bytes)
}

func (k *KoreanBasic) getSearch(q string) ([]byte, error) {
	resp, err := http.Get(apiURL(q, k.apiKey))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response from API is not OK (200), got: %v (%v)", resp.Status, resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

// SearchTerms returns the search terms given a search response
func SearchTerms(searchResponse []byte) ([]dictionary.Term, error) {
	channel, err := unmarshallSearch(searchResponse)
	if err != nil {
		return nil, err
	}
	return itemsToTerms(channel.Items)
}

const apiURLString = "https://krdict.korean.go.kr/api/search?sort=popular&translated=y&trans_lang=1&q=%s&key=%s"

func apiURL(q, apiKey string) string {
	return fmt.Sprintf(apiURLString, url.QueryEscape(q), apiKey)
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
	"명사":     lang.PartOfSpeechNoun,
	"대명사":    lang.PartOfSpeechPronoun,
	"수사":     lang.PartOfSpeechNumeral,
	"조사":     lang.PartOfSpeechPostposition,
	"동사":     lang.PartOfSpeechVerb,
	"형용사":    lang.PartOfSpeechAdjective,
	"관형사":    lang.PartOfSpeechPrenoun,
	"부사":     lang.PartOfSpeechAdverb,
	"감탄사":    lang.PartOfSpeechInterjection,
	"접사":     lang.PartOfSpeechAffix,
	"의존 명사":  lang.PartOfSpeechDependentNoun,
	"보조 동사":  lang.PartOfSpeechAuxiliaryVerb,
	"보조 형용사": lang.PartOfSpeechAuxiliaryAdjective,
	"어미":     lang.PartOfSpeechEnding,
	"품사 없음":  lang.PartOfSpeechNone,
}

func itemsToTerms(items []item) ([]dictionary.Term, error) {
	terms := make([]dictionary.Term, len(items))
	for i, item := range items {
		if _, exists := partOfSpeechMap[item.PartOfSpeech]; !exists {
			return nil, fmt.Errorf("part of speech not found: %v, %v", item.Word, item.PartOfSpeech)
		}
		terms[i] = dictionary.Term{
			Text:             item.Word,
			CommonLevel:      wordGradeToCommonLevel[item.WordGrade],
			PartOfSpeech:     partOfSpeechMap[item.PartOfSpeech],
			DictionarySource: DictionarySource,
		}
		terms[i].Translations = make([]dictionary.Translation, len(item.Senses))
		for j, sense := range item.Senses {
			terms[i].Translations[j] = dictionary.Translation{
				Text:        sense.Translation,
				Explanation: sense.Explanation,
			}
		}
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
	Word         string  `xml:"word"`
	WordGrade    string  `xml:"word_grade"`
	PartOfSpeech string  `xml:"pos"`
	Senses       []sense `xml:"sense"`
}

type sense struct {
	Translation string `xml:"translation>trans_word"`
	Explanation string `xml:"translation>trans_dfn"`
}
