package khaiii

/*
#include <string.h>

#define RESERVED_STRLEN 8

typedef struct khaiii_morph_t_ {
    const char* lex;    ///< lexical
    const char* tag;    ///< part-of-speech tag
    int begin;    ///< morpheme begin position
    int length;    ///< morpheme length
    char reserved[RESERVED_STRLEN];    ///< reserved
    const struct khaiii_morph_t_* next;    ///< next pointer
} khaiii_morph_t;

typedef struct khaiii_word_t_ {
    int begin;    ///< word begin position
    int length;    ///< word length
    char reserved[RESERVED_STRLEN];    ///< reserved
    const khaiii_morph_t* morphs;    ///< morpheme list
    const struct khaiii_word_t_* next;    ///< next pointer
} khaiii_word_t;
*/
import "C"

import (
	"fmt"
	"unsafe"
)

var version func() string
var open func(rscDir, optStr string) int
var analyze func(handle int, input, optStr string) *C.khaiii_word_t
var freeResults func(handle int, results *C.khaiii_word_t)
var close func(handle int)
var lastError func(handle int) string

var apiFunctions = map[string]interface{}{
	"khaiii_version":      &version,
	"khaiii_open":         &open,
	"khaiii_analyze":      &analyze,
	"khaiii_free_results": &freeResults,
	"khaiii_close":        &close,
	"khaiii_last_error":   &lastError,
}

// Morph represents a morpheme in a word
type Morph struct {
	Lex      string
	Tag      string
	Begin    int
	Length   int
	reserved string
}

// Word represents a word in the given string
type Word struct {
	Begin    int
	Length   int
	reserved string
	Morphs   []*Morph
}

// Version returns the version of Khaiii being run
func (k *Khaiii) Version() string {
	return version()
}

// Open opens the training resource directory
func (k *Khaiii) Open(rscDir string) error {
	if k.openHandle != 0 {
		return fmt.Errorf("Khaiii.Open() is already open")
	}
	openHandle := open(rscDir, "{}")
	if openHandle == -1 {
		return fmt.Errorf("failed to Khaiii.Open(): %s", k.lastError())
	}
	k.openHandle = openHandle
	return nil
}

// Analyze analyzes the input string
func (k *Khaiii) Analyze(input string) ([]*Word, error) {
	if k.openHandle <= 0 {
		return nil, fmt.Errorf("Khaiii.Open() invalid for Analyze()")
	}

	var err error
	wordC := analyze(k.openHandle, input, "{}")
	if wordC == nil {
		return nil, fmt.Errorf("failed to Khaiii.Analyze(): %s", k.lastError())
	}
	defer func() {
		err2 := k.freeResults(wordC)
		if err == nil {
			err = err2
		}
	}()

	words := []*Word{}
	for wordC != nil {
		word := goWord(*wordC)
		words = append(words, word)
		wordC = wordC.next
	}

	return words, err
}

// Close cleans Khaiii's resources, so it can stop
func (k *Khaiii) Close() error {
	if k.openHandle <= 0 {
		return fmt.Errorf("trying to Khaiii.Close() with handle: %v", k.openHandle)
	}
	close(k.openHandle)
	k.openHandle = 0
	return nil
}

func (k *Khaiii) freeResults(wordC *C.khaiii_word_t) error {
	freeResults(k.openHandle, wordC)
	return nil
}

func (k *Khaiii) lastError() string {
	return lastError(k.openHandle)
}

func goWord(wordC C.khaiii_word_t) *Word {
	word := &Word{
		Begin:    int(wordC.begin),
		Length:   int(wordC.length),
		reserved: strndup((*C.char)(unsafe.Pointer(&wordC.reserved)), C.RESERVED_STRLEN),
	}

	morphs := []*Morph{}
	morphC := wordC.morphs
	for morphC != nil {
		morph := goMorph(*morphC)
		morphs = append(morphs, morph)
		morphC = morphC.next
	}
	word.Morphs = morphs
	return word
}

func goMorph(morphC C.khaiii_morph_t) *Morph {
	morph := &Morph{
		Lex:      C.GoString(morphC.lex),
		Tag:      C.GoString(morphC.tag),
		Begin:    int(morphC.begin),
		Length:   int(morphC.length),
		reserved: strndup((*C.char)(unsafe.Pointer(&morphC.reserved)), C.RESERVED_STRLEN),
	}
	return morph
}

func strndup(cs *C.char, len int) string {
	return C.GoStringN(cs, C.int(C.strnlen(cs, C.size_t(len))))
}
