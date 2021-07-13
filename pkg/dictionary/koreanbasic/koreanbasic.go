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
	return itemsToTerms(channel.Items), nil
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

func itemsToTerms(items []item) []dictionary.Term {
	terms := make([]dictionary.Term, len(items))
	for i, item := range items {
		terms[i] = dictionary.Term{
			Text:        item.Word,
			CommonLevel: wordGradeToCommonLevel[item.WordGrade],
		}
		terms[i].Translations = make([]dictionary.Translation, len(item.Senses))
		for j, sense := range item.Senses {
			terms[i].Translations[j] = dictionary.Translation{
				Text:        sense.Translation,
				Explanation: sense.Explanation,
			}
		}
	}
	return terms
}

type channel struct {
	Items  []item `xml:"item"`
	Total  uint   `xml:"total"`
	Start  uint   `xml:"start"`
	Number uint   `xml:"num"`
}

type item struct {
	Word      string  `xml:"word"`
	WordGrade string  `xml:"word_grade"`
	Senses    []sense `xml:"sense"`
}

type sense struct {
	Translation string `xml:"translation>trans_word"`
	Explanation string `xml:"translation>trans_dfn"`
}
