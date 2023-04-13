package krdict

import (
	"encoding/xml"
	"os"
)

// Seed seeds the database from the rscPath XML
func Seed() error {
	return nil
}

func unmarshallRscPath() ([]*lexicalResource, error) {
	lexes := []*lexicalResource{}
	xmlPaths, err := RscXMLPaths()
	if err != nil {
		return nil, err
	}
	for _, path := range xmlPaths {
		//nolint:gosec // just parsing XML
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		lex, err := unmarshallXML(bytes)
		if err != nil {
			return nil, err
		}
		lexes = append(lexes, lex)
	}
	return lexes, nil
}

func unmarshallXML(bytes []byte) (*lexicalResource, error) {
	lex := &lexicalResource{}
	if err := xml.Unmarshal(bytes, lex); err != nil {
		return nil, err
	}
	return lex, nil
}

func findGoodExample() (*lexicalEntry, error) {
	lexes, err := unmarshallRscPath()
	if err != nil {
		return nil, err
	}
	for _, lex := range lexes {
		for _, entry := range lex.LexicalEntries {
			if goodExampleEntry(entry) && goodExampleSense(entry.Senses) && goodExampleWordForm(entry.WordForms) {
				return &entry, nil
			}
		}
	}
	return nil, nil
}

func goodExampleEntry(entry lexicalEntry) bool {
	return !(entry.Lemmas == nil || entry.RelatedForms == nil || entry.Senses == nil || entry.WordForms == nil)
}

func goodExampleSense(senses []sense) bool {
	for _, sense := range senses {
		if !(sense.Equivalents == nil || sense.Multimedias == nil || sense.SenseExamples == nil || sense.SenseRelations == nil) {
			return true
		}
	}
	return false
}

func goodExampleWordForm(wordForms []wordForm) bool {
	for _, wordForm := range wordForms {
		if !(wordForm.FormRepresentation.Feats == nil) {
			return true
		}
	}
	return false
}

type lexicalResource struct {
	LexicalEntries []lexicalEntry `xml:"Lexicon>LexicalEntry"`
}

type lexicalEntry struct {
	ID    uint   `xml:"val,attr"`
	Feats []feat `xml:"feat"`

	Lemmas       []lemma       `xml:"Lemma"`
	RelatedForms []relatedForm `xml:"RelatedForm"`
	Senses       []sense       `xml:"Sense"`
	WordForms    []wordForm    `xml:"WordForm"`
}

type lemma struct {
	Feats []feat `xml:"feat"`
}

type relatedForm struct {
	Feats []feat `xml:"feat"`
}

type sense struct {
	ID    uint   `xml:"val,attr"`
	Feats []feat `xml:"feat"`

	Equivalents    []equivalent    `xml:"Equivalent"`
	Multimedias    []multimedia    `xml:"Multimedia"`
	SenseExamples  []senseExample  `xml:"SenseExample"`
	SenseRelations []senseRelation `xml:"SenseRelation"`
}

type equivalent struct {
	Feats []feat `xml:"feat"`
}

type multimedia struct {
	Feats []feat `xml:"feat"`
}

type senseExample struct {
	Feats []feat `xml:"feat"`
}

type senseRelation struct {
	Feats []feat `xml:"feat"`
}

type wordForm struct {
	Feats              []feat             `xml:"feat"`
	FormRepresentation formRepresentation `xml:"FormRepresentation"`
}

type formRepresentation struct {
	Feats []feat `xml:"feat"`
}

type feat struct {
	Att string `xml:"att,attr"`
	Val string `xml:"val,attr"`
}
