// Package khaiii is an interface to the Khaiii Korean tokenizer
package khaiii

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/server"
	api "github.com/s12chung/text2anki/tokenizers/khaiii/pkg/khaiii"
)

const stopWarningDuration = 15 * time.Second
const port = 9999
const binName = "./khaiiiserver"

var binPath string

func init() {
	binPath = os.Getenv("KHAIII_BIN_PATH")
}

// New returns a Khaiii Korean tokenizer
func New() tokenizers.Tokenizer {
	return new()
}

func new() *Khaiii {
	name := "Khaiii"
	server := server.NewCmdTokenizerServer(port, stopWarningDuration,
		binPath,
		binName, "--port", strconv.Itoa(port))
	return &Khaiii{
		name:            name,
		ServerTokenizer: tokenizers.NewServerTokenizer(name, server),
	}
}

// Khaiii is a Korean Tokenizer in java
type Khaiii struct {
	name string
	tokenizers.ServerTokenizer
}

// Tokenize returns the part of speech tokens of the given string
func (k *Khaiii) Tokenize(str string) ([]tokenizers.Token, error) {
	resp := &api.TokenizeResponse{}
	err := k.ServerTokenize(str, resp)
	if err != nil {
		return nil, err
	}
	return k.toTokenizerTokens(resp)
}

func (k *Khaiii) toTokenizerTokens(resp *api.TokenizeResponse) ([]tokenizers.Token, error) {
	tokens := []tokenizers.Token{}
	for _, word := range resp.Words {
		for _, morph := range word.Morphs {
			partOfSpeech, found := partOfSpeechMap[morph.Tag]
			if !found {
				return nil, fmt.Errorf("%v POS not mapped: %v", k.name, morph.Tag)
			}
			tokens = append(tokens, tokenizers.Token{
				Text:         morph.Lex,
				PartOfSpeech: partOfSpeech,
				StartIndex:   morph.Begin,
				EndIndex:     morph.Begin + morph.Length,
			})
		}
	}
	return tokens, nil
}

var partOfSpeechMap = map[string]lang.PartOfSpeech{
	"NNG": lang.PartOfSpeechNoun,          // Common Noun - 체언 - 일반 명사
	"NNP": lang.PartOfSpeechNoun,          // Proper Noun/Name - 체언 - 고유 명사
	"NNB": lang.PartOfSpeechDependentNoun, // 체언 - 의존 명사
	"NP":  lang.PartOfSpeechPronoun,       // 체언 - 대명사
	"NR":  lang.PartOfSpeechNumeral,       // 체언 - 수사

	"VV":  lang.PartOfSpeechVerb,               // 용언 - 동사
	"VA":  lang.PartOfSpeechAdjective,          // 용언 - 형용사
	"VX":  lang.PartOfSpeechAuxiliaryPredicate, // 용언 - 보조 용언
	"VCP": lang.PartOfSpeechCopula,             // Positive Copula - 용언 - 긍정 지정사
	"VCN": lang.PartOfSpeechCopula,             // Negative Copula - 용언 - 부정 지정사

	"MM":  lang.PartOfSpeechDeterminer,   // 수식언 - 관형사
	"MAG": lang.PartOfSpeechAdverb,       // General Adverb - 수식언 - 일반 부사
	"MAJ": lang.PartOfSpeechAdverb,       // Joining Adverb - 수식언 - 접속 부사
	"IC":  lang.PartOfSpeechInterjection, // 독립언 - 감탄사

	"JKS": lang.PartOfSpeechPostposition, // Subject Case Marker - 관계언 - 주격 조사
	"JKC": lang.PartOfSpeechPostposition, // Complement Case Marker - 관계언 - 보격 조사
	"JKG": lang.PartOfSpeechPostposition, // Adnominal Case Marker - 관계언 - 관형격 조사
	"JKO": lang.PartOfSpeechPostposition, // Object Case Marker - 관계언 - 목적격 조사
	"JKB": lang.PartOfSpeechPostposition, // Adverbial Case Marker - 관계언 - 부사격 조사
	"JKV": lang.PartOfSpeechPostposition, // Vocative Case Marker - 관계언 - 호격 조사
	"JKQ": lang.PartOfSpeechPostposition, // Quoted Case Particle - 관계언 - 인용격 조사
	"JX":  lang.PartOfSpeechPostposition, // Auxiliary Postpositional Particle - 관계언 - 보조사
	"JC":  lang.PartOfSpeechPostposition, // Conjunctive Postpositional Particle - 관계언 - 접속 조사

	"EP":  lang.PartOfSpeechEnding, // Pre-Final Ending - 의존 형태 - 선어말 어미
	"EF":  lang.PartOfSpeechEnding, // Sentence Closing Ending - 의존 형태 - 종결 어미
	"EC":  lang.PartOfSpeechEnding, // Connective Ending - 의존 형태 - 연결 어미
	"ETN": lang.PartOfSpeechEnding, // Nominal Ending - 의존 형태 - 명사형 전성 어미
	"ETM": lang.PartOfSpeechEnding, // Transformative Ending - 의존 형태 - 관형형 전성 어미

	"XPN": lang.PartOfSpeechPrefix, // Noun Prefix - 의존 형태 - 체언 접두사
	"XSN": lang.PartOfSpeechSuffix, // Noun Derived Suffix - 의존 형태 - 명사 파생 접미사
	"XSV": lang.PartOfSpeechSuffix, // Verb Derived Suffix - 의존 형태 - 동사 파생 접미사
	"XSA": lang.PartOfSpeechSuffix, // Adjective Derived Suffix - 의존 형태 - 형용사 파생 접미사

	"XR": lang.PartOfSpeechRoot, // Root - 의존 형태 - 어근

	"SF": lang.PartOfSpeechPunctuation, // Period, Question Mark, Exclamation Mark - 기호 - 마침표, 물음표, 느낌표
	"SP": lang.PartOfSpeechPunctuation, // Comma, Bullet, Colon, Slash - 기호 - 쉼표, 가운뎃점, 콜론, 빗금
	"SS": lang.PartOfSpeechPunctuation, // Quotation Mark, Parentheses, Dash - 기호 - 따옴표, 괄호표, 줄표
	"SE": lang.PartOfSpeechPunctuation, // Ellipsis - 기호 - 줄임표
	"SO": lang.PartOfSpeechPunctuation, // Dash, Tilde, Hidden - 기호 - 붙임표 (물결, 숨김, 빠짐)
	"SW": lang.PartOfSpeechPunctuation, // Other Symbols (Logic, Mathematical, Monetary, etc.) - 기호 - 기타 기호 (논리, 수학 기호, 화폐 기호 등)

	"SL": lang.PartOfSpeechOtherLanguage, // Foreign Language - 기호 - 외국어
	"SH": lang.PartOfSpeechOtherLanguage, // Chinese Language - 기호 - 한자
	"SN": lang.PartOfSpeechNumeral,       // 기호 - 숫자

	// All below are not in Sejeong Corpus
	"SWK": lang.PartOfSpeechAlphabet, // Korean Alphabet (subpart of SW) - 기호 - 한글 자소

	"ZN": lang.PartOfSpeechNoun,  // Guessed Noun  - 추정 - 분석 불능 (명사 추정)
	"ZV": lang.PartOfSpeechVerb,  // Guessed Verb - 추정 - 분석 불능 (용언 추정)
	"ZZ": lang.PartOfSpeechOther, // Guessed Other - 추정 - 분석 불능 (기타)
}
