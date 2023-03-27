// Package khaiii is an interface to the Khaiii Korean tokenizer
package khaiii

import "github.com/s12chung/text2anki/pkg/lang"

var partOfSpeechMap = map[string]lang.PartOfSpeech{
	// NNG	common noun	체언	일반 명사
	// NNP	proper noun	체언	고유 명사
	// NNB	dependent noun	체언	의존 명사
	// NP	pronoun	체언	대명사
	// NR	numeral	체언	수사
	// VV	Verb	용언	동사
	// VA	an adjective	용언	형용사
	// VX	supplementary idiom	용언	보조 용언
	// VCP	affirmative specifier	용언	긍정 지정사
	// VCN	indefinite designation	용언	부정 지정사
	// MM	tubular detective	수식언	관형사
	// MAG	general adverb	수식언	일반 부사
	// MAJ	conjunction adverb	수식언	접속 부사
	// IC	Interjection	독립언	감탄사
	// JKS	subjective investigation	관계언	주격 조사
	// JKC	complementary investigation	관계언	보격 조사
	// JKG	tubular case study	관계언	관형격 조사
	// JKO	objective case study	관계언	목적격 조사
	// JKB	adverbial investigation	관계언	부사격 조사
	// JKV	favoritism	관계언	호격 조사
	// JKQ	citation investigation	관계언	인용격 조사
	// JX	auxiliary verb	관계언	보조사
	// JC	connection survey	관계언	접속 조사
	// EP	predicate ending	의존 형태	선어말 어미
	// EF	terminating ending	의존 형태	종결 어미
	// EC	concatenated ending	의존 형태	연결 어미
	// ETN	noun-type prime ending	의존 형태	명사형 전성 어미
	// ETM	tubular transitive ending	의존 형태	관형형 전성 어미
	// XPN	phonetic prefix	의존 형태	체언 접두사
	// XSN	noun-derived suffix	의존 형태	명사 파생 접미사
	// XSV	verb-derived suffix	의존 형태	동사 파생 접미사
	// XSA	adjective-derived suffix	의존 형태	형용사 파생 접미사
	// XR	root of a word	의존 형태	어근
	// SF	Period, question mark, exclamation mark.	기호	마침표, 물음표, 느낌표
	// SP	comma, middle point, colon, and comb.	기호	쉼표, 가운뎃점, 콜론, 빗금
	// SS	Quotes, parentheses, and lines.	기호	따옴표, 괄호표, 줄표
	// SE	abbreviation table	기호	줄임표
	// SO	Attachment mark (wave, hidden, missing)	기호	붙임표(물결, 숨김, 빠짐)
	// SL	Foreign language	기호	외국어
	// SH	Chinese characters	기호	한자
	// SW	Other symbols (logic, mathematical, monetary, etc.)	기호	기타 기호(논리, 수학 기호, 화폐 기호 등)
	// SWK	Korean alphabet	기호	한글 자소
	// SN	number	기호	숫자
	// ZN	Non-analytic (noun estimation)	추정	분석 불능(명사 추정)
	// ZV	Non-analytic (prediction)	추정	분석 불능(용언 추정)
	// ZZ	non-analytic (other)	추정	분석 불능(기타)
}
