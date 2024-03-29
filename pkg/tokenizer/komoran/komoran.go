// Package komoran is an interface to the Komoran Korean tokenizer
package komoran

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/tokenizer"
	"github.com/s12chung/text2anki/pkg/tokenizer/server"
)

const stopWarningDuration = 15 * time.Second
const port = 9999
const backlog = 64
const jarName = "tokenizer-komoran.jar"

var jarPath string

func init() {
	jarPath = os.Getenv("KOMORAN_JAR_PATH")
}

// New returns a Komoran Korean tokenizer
func New(ctx context.Context, log *slog.Logger) tokenizer.Tokenizer {
	return newKomoran(ctx, port, log)
}

func newKomoran(ctx context.Context, port int, log *slog.Logger) *Komoran {
	opts := server.NewCmdOptions("java")
	opts.Args = []string{"-jar", jarName, strconv.Itoa(port), strconv.Itoa(backlog)}
	opts.Dir = jarPath
	server := server.NewCmdTokenizerServer(ctx, opts, port, stopWarningDuration, log)

	name := "Komoran"
	return &Komoran{
		name:            name,
		ServerTokenizer: tokenizer.NewServerTokenizer(name, server),
	}
}

// Komoran is a Korean Tokenizer in java
type Komoran struct {
	name string
	tokenizer.ServerTokenizer
}
type response struct {
	Tokens []token `json:"tokens"`
}

// nolint:tagliatelle // following java format
type token struct {
	POS        string `json:"pos"`
	EndIndex   uint   `json:"endIndex"`
	BeginIndex uint   `json:"beginIndex"`
	Morph      string `json:"morph"`
}

// Tokenize returns the part of speech tokens of the given string
func (k *Komoran) Tokenize(ctx context.Context, str string) ([]tokenizer.Token, error) {
	resp := &response{}
	err := k.ServerTokenize(ctx, str, resp)
	if err != nil {
		return nil, err
	}
	return k.toTokenizerTokens(resp)
}

func (k *Komoran) toTokenizerTokens(resp *response) ([]tokenizer.Token, error) {
	tokens := make([]tokenizer.Token, len(resp.Tokens))
	for i, token := range resp.Tokens {
		partOfSpeech, found := partOfSpeechMap[token.POS]
		if !found {
			return nil, fmt.Errorf("%v POS not mapped: %v", k.name, token.POS)
		}
		tokens[i] = tokenizer.Token{
			Text:         token.Morph,
			PartOfSpeech: partOfSpeech,
			StartIndex:   token.BeginIndex,
			Length:       token.EndIndex - token.BeginIndex,
		}
	}
	return tokens, nil
}

// nolint:lll // link is just long
//
// There are symbols that do not exist in running source code returned. See Google Doc
//
// https://docs.google.com/spreadsheets/d/1OGAjUvalBuX-oZvZ_-9tEfYD2gQe7hTGsgUpiiBSXI8/edit#gid=0
// https://github.com/shineware/KOMORAN/blob/d0badf1f947c6ba60f5cb16b5b1e5fa61b69aad2/core/src/main/java/kr/co/shineware/nlp/komoran/constant/SYMBOL.java

var partOfSpeechMap = map[string]lang.PartOfSpeech{
	"NNG": lang.PartOfSpeechNoun,    // Common Noun - 일반 명사
	"NNP": lang.PartOfSpeechNoun,    // Proper Noun/Name - 고유 명사
	"NP":  lang.PartOfSpeechPronoun, // 대명사
	"NR":  lang.PartOfSpeechNumeral, // 수사
	"SN":  lang.PartOfSpeechNumeral, // 숫자 - Not in Github

	"JC":  lang.PartOfSpeechPostposition, // Conjunctive Postpositional Particle - 접속 조사
	"JKB": lang.PartOfSpeechPostposition, // Adverbial Case Marker - 부사격 조사
	"JKC": lang.PartOfSpeechPostposition, // Complement Case Marker - 보격 조사
	"JKG": lang.PartOfSpeechPostposition, // Adnominal Case Marker - 관형격 조사
	"JKO": lang.PartOfSpeechPostposition, // Object Case Marker - 목적격 조사
	"JKS": lang.PartOfSpeechPostposition, // Subject Case Marker - 주격 조사
	"JKV": lang.PartOfSpeechPostposition, // Vocative Case Marker 호격 조사
	"JX":  lang.PartOfSpeechPostposition, // Auxiliary Postpositional Particle - 보조사
	"JKQ": lang.PartOfSpeechPostposition, // Quoted Case Particle - 인용격 조사 - Not in Github

	"VV":  lang.PartOfSpeechVerb,         // 동사
	"VA":  lang.PartOfSpeechAdjective,    // 형용사
	"MM":  lang.PartOfSpeechDeterminer,   // 관형사 - Not in Github
	"MAG": lang.PartOfSpeechAdverb,       // General Adverb - 일반 부사 - Not in Github
	"MAJ": lang.PartOfSpeechAdverb,       // Joining Adverb - 접속 부사 - Not in Github
	"IC":  lang.PartOfSpeechInterjection, // 감탄사 - Not in Github

	"XPN": lang.PartOfSpeechPrefix, // Noun Prefix - 체언 접두사 - Not in GitHub
	"XSN": lang.PartOfSpeechSuffix, // Noun Derived Suffix - 명사 파생 접미사 - Not in Github
	"XSV": lang.PartOfSpeechSuffix, // Verb Derived Suffix - 동사 파생 접미사 - Not in Github
	"XSA": lang.PartOfSpeechSuffix, // Adjective Derived Suffix - 형용사 파생 접미사 - Not in Github

	"NNB": lang.PartOfSpeechDependentNoun,      // 의존 명사
	"VX":  lang.PartOfSpeechAuxiliaryPredicate, // 보조 용언

	"EP":  lang.PartOfSpeechEnding, // Pre-Final Ending 선어말 어미
	"EC":  lang.PartOfSpeechEnding, // Connective Ending - 연결 어미
	"EF":  lang.PartOfSpeechEnding, // Sentence Closing Ending - 종결 어미
	"ETN": lang.PartOfSpeechEnding, // Nominal Ending - 명사형 전성 어미
	"ETM": lang.PartOfSpeechEnding, // Transformative Ending - 관형형 전성 어미
	"VCP": lang.PartOfSpeechCopula, // Positive Copula - 긍정 지정사
	"VCN": lang.PartOfSpeechCopula, // Negative Copula - 부정 지정사

	"SW": lang.PartOfSpeechPunctuation, // Other Symbols - 기타 기호
	"SF": lang.PartOfSpeechPunctuation, // Period, Question Mark, Exclamation Mark - 마침표, 물음표, 느낌표
	"SS": lang.PartOfSpeechPunctuation, // Quotation Mark, Parentheses, Dash - 따옴표, 괄호표, 줄표
	"SE": lang.PartOfSpeechPunctuation, // Ellipsis - 줄임표 - Not in Github
	"SP": lang.PartOfSpeechPunctuation, // Comma, Bullet, Colon, Slash - 쉼표, 가운뎃점, 콜론, 빗금 - Not in Github
	"SO": lang.PartOfSpeechPunctuation, // Dash, Tilde, Hidden - 붙임표 (물결, 숨김, 빠짐) - Not in Github

	"SH": lang.PartOfSpeechOtherLanguage, // Chinese Language - 한자 - Not in Github
	"SL": lang.PartOfSpeechOtherLanguage, // Foreign Language - 외국어 - Not in Github
	"XR": lang.PartOfSpeechRoot,          // Root - 어근 - Not in Github

	"NA": lang.PartOfSpeechUnknown, // Unknown - 분석불능범주
	"NF": lang.PartOfSpeechUnknown, // Presumptive Noun Category of Nouns - 명사추정범주 - Not in Github
	"NV": lang.PartOfSpeechUnknown, // Prediction Category of Terminology - 용언추정범주 - Not in Github
}
