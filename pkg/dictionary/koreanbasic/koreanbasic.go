// Package koreanbasic contains functions for the Korean Basic Dictionary API
package koreanbasic

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/s12chung/text2anki/pkg/dictionary"
)

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
	resp, err := http.Get(apiURL(q, k.apiKey))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response from API is not OK (200), got: %v (%v)", resp.Status, resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	channel, err := parseSearch(bytes)
	if err != nil {
		return nil, err
	}
	return itemsToTerms(channel.Items)
}

const apiURLString = "https://krdict.korean.go.kr/api/search?sort=popular&translated=y&trans_lang=1&q=%s&key=%s"

func apiURL(q, apiKey string) string {
	return fmt.Sprintf(apiURLString, url.QueryEscape(q), apiKey)
}

func parseSearch(data []byte) (*channel, error) {
	ch := &channel{}
	if err := xml.Unmarshal(data, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

var wordGradeToCommonLevel = map[string]dictionary.CommonLevel{
	"":   dictionary.CommonLevelUnique,
	"고급": dictionary.CommonLevelRare,
	"중급": dictionary.CommonLevelMedium,
	"초급": dictionary.CommonLevelCommon,
}

var partOfSpeechMap = map[string]dictionary.PartOfSpeech{
	"명사":     dictionary.PartOfSpeechNoun,
	"대명사":    dictionary.PartOfSpeechPronoun,
	"수사":     dictionary.PartOfSpeechNumeral,
	"조사":     dictionary.PartOfSpeechPostposition,
	"동사":     dictionary.PartOfSpeechVerb,
	"형용사":    dictionary.PartOfSpeechAdjective,
	"관형사":    dictionary.PartOfSpeechPrenoun,
	"부사":     dictionary.PartOfSpeechAdverb,
	"감탄사":    dictionary.PartOfSpeechInterjection,
	"접사":     dictionary.PartOfSpeechAffix,
	"의존 명사":  dictionary.PartOfSpeechDependentNoun,
	"보조 동사":  dictionary.PartOfSpeechAuxiliaryVerb,
	"보조 형용사": dictionary.PartOfSpeechAuxiliaryAdjective,
	"어미":     dictionary.PartOfSpeechEnding,
	"품사 없음":  dictionary.PartOfSpeechNone,
}

func itemsToTerms(items []item) ([]dictionary.Term, error) {
	terms := make([]dictionary.Term, len(items))
	for i, item := range items {
		if _, exists := partOfSpeechMap[item.PartOfSpeech]; !exists {
			return nil, fmt.Errorf("part of speech not found: %v, %v", item.Word, item.PartOfSpeech)
		}
		terms[i] = dictionary.Term{
			Text:         item.Word,
			CommonLevel:  wordGradeToCommonLevel[item.WordGrade],
			PartOfSpeech: partOfSpeechMap[item.PartOfSpeech],
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
