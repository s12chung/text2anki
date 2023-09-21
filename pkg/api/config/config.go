// Package config contains the config for package api
package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/dictionary/krdict"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/extractor/instagram"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/synthesizer"
	"github.com/s12chung/text2anki/pkg/synthesizer/azure"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizer"
	"github.com/s12chung/text2anki/pkg/tokenizer/khaiii"
	"github.com/s12chung/text2anki/pkg/tokenizer/komoran"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx"
	"github.com/s12chung/text2anki/pkg/util/logg"
)

var appCacheDir string

func init() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		slog.Error("config.init()", logg.Err(err)) //nolint:forbidigo // used in init only
		os.Exit(-1)
	}
	appCacheDir = path.Join(cacheDir, "Text2Anki")
}

// Config contains config settings for the API
type Config struct {
	CacheDir string
	Log      *slog.Logger

	UUIDGenerator storage.UUIDGenerator
	TxPool        reqtx.Pool[db.TxQs, TxMode]
	StorageConfig StorageConfig

	TokenizerType
	DictionaryType

	ExtractorMap extractor.Map

	Synthesizer  synthesizer.Synthesizer
	AnkiCacheDir string
}

// CacheDir returns the cache dir for the API
func CacheDir(cacheDir string) string {
	if cacheDir == "" {
		cacheDir = appCacheDir
	}
	return cacheDir
}

// UUIDGenerator returns the configured uuidGenerator
func UUIDGenerator(uuidGenerator storage.UUIDGenerator) storage.UUIDGenerator {
	if uuidGenerator == nil {
		uuidGenerator = storage.UUID7{}
	}
	return uuidGenerator
}

// TxMode identifies the transaction mode
type TxMode int

// TxIntegrator returns a new TxIntegrator
func TxIntegrator(txPool reqtx.Pool[db.TxQs, TxMode]) reqtx.Integrator[db.TxQs, TxMode] {
	return reqtx.NewIntegrator(txPool)
}

// StorageType defines what signer type to for file storage
type StorageType int

const (
	// StorageLocalStore picks the local file store
	StorageLocalStore StorageType = iota
)

// StorageConfig configures the storage signer
type StorageConfig struct {
	StorageType
	LocalStoreConfig LocalStoreConfig
	UUIDGenerator    storage.UUIDGenerator
}

// Storage contains the Route's storage setup
type Storage struct {
	DBStorage storage.DBStorage
	Storer    storage.Storer
}

// StorageFromConfig returns a storage from the given config
func StorageFromConfig(config StorageConfig, log *slog.Logger) Storage {
	var storageAPI storage.API
	var storer storage.Storer
	var err error
	switch config.StorageType {
	case StorageLocalStore:
		fallthrough
	default:
		var ls localstore.API
		ls, err = LocalStoreAPI(config.LocalStoreConfig)
		storageAPI = ls
		storer = ls
	}
	if err != nil {
		log.Error("config.StorageFromConfig()", logg.Err(err))
		os.Exit(-1)
	}
	return Storage{DBStorage: storage.NewDBStorage(storageAPI, config.UUIDGenerator), Storer: storer}
}

// LocalStoreConfig defines the config for localstore
type LocalStoreConfig struct {
	Origin        string
	KeyBasePath   string
	EncryptorPath string
}

var localStoreConfigValidator = firm.NewStructValidator(firm.RuleMap{
	"Origin":        {rule.Present{}},
	"KeyBasePath":   {rule.Present{}},
	"EncryptorPath": {rule.Present{}},
})

const localstoreKey = "localstore.key"

// StorageURLPath is the default storage URL path for the LocalStoreAPI
const StorageURLPath = "/storage"

// LocalStoreAPI returns a localstore.API
func LocalStoreAPI(config LocalStoreConfig) (localstore.API, error) {
	if !strings.HasSuffix(config.Origin, "/") {
		config.Origin += "/"
	}
	config.Origin += StorageURLPath[1:]

	// LocalStoreAPI is called when declaring package level vars (before init()), this ensures the definition works
	result := localStoreConfigValidator.Validate(config)
	if !result.IsValid() {
		return localstore.API{}, fmt.Errorf(result.ErrorMap().String())
	}
	encryptor, err := localstore.NewAESEncryptorFromFile(path.Join(config.EncryptorPath, localstoreKey))
	if err != nil {
		return localstore.API{}, err
	}
	return localstore.NewAPI(config.Origin, config.KeyBasePath, encryptor), nil
}

// Parser returns the default Parser
func Parser() text.Parser { return text.NewParser(text.Korean, text.English) }

// TokenizerType is an enum of tokenizer types
type TokenizerType int

const (
	// TokenizerKhaiii picks the Khaiii tokenizer
	TokenizerKhaiii TokenizerType = iota
	// TokenizerKomoran picks the Komoran tokenizer
	TokenizerKomoran
)

// Tokenizer returns the default Tokenizer
func Tokenizer(ctx context.Context, tokenizerType TokenizerType, log *slog.Logger) tokenizer.Tokenizer {
	switch tokenizerType {
	case TokenizerKomoran:
		return komoran.New(ctx, log)
	case TokenizerKhaiii:
		fallthrough
	default:
		return khaiii.New(ctx, log)
	}
}

// DictionaryType is an enum of dictionary types
type DictionaryType int

const (
	// DictionaryKrDict picks the KrDict dictionary
	DictionaryKrDict DictionaryType = iota
	// DictionaryKoreanBasic picks the KoreanBasic dictionary
	DictionaryKoreanBasic
)

// Dictionary returns the default Dictionary
func Dictionary(dictionaryType DictionaryType) dictionary.Dictionary {
	switch dictionaryType {
	case DictionaryKoreanBasic:
		return koreanbasic.New(koreanbasic.GetAPIKeyFromEnv())
	case DictionaryKrDict:
		fallthrough
	default:
		return krdict.New()
	}
}

// ExtractorMap returns the ExtractorMap config
func ExtractorMap(extractorMap extractor.Map) extractor.Map {
	if extractorMap != nil {
		return extractorMap
	}
	return extractor.Map{
		"instagram": extractor.NewExtractor(filepath.Join(appCacheDir, "instagram"), instagram.Factory{}),
	}
}

// Synthesizer returns the default Synthesizer
func Synthesizer(synth synthesizer.Synthesizer) synthesizer.Synthesizer {
	if synth == nil {
		synth = azure.New(azure.GetAPIKeyFromEnv(), azure.EastUSRegion)
	}
	return synth
}

// SynthesizerSoundFactory is the SoundFactory for the Synthesizer
type SynthesizerSoundFactory struct {
	synth synthesizer.Synthesizer
}

// Name returns the name of the Synthesizer
func (s SynthesizerSoundFactory) Name() string { return s.synth.SourceName() }

// Sound returns the sound from the Synthesizer
func (s SynthesizerSoundFactory) Sound(ctx context.Context, usage string) ([]byte, error) {
	return s.synth.TextToSpeech(ctx, usage)
}

// SoundSetter returns the SoundSetter given the Synthesizer
func SoundSetter(synth synthesizer.Synthesizer) anki.SoundSetter {
	return anki.NewSoundSetter(SynthesizerSoundFactory{synth: synth})
}

// Anki returns the anki.Config
func Anki(cacheDir string, log *slog.Logger) anki.Config {
	if cacheDir == "" {
		cacheDir = path.Join(appCacheDir, "notes")
	}
	if err := os.MkdirAll(cacheDir, ioutil.OwnerRWXGroupRX); err != nil {
		log.Error("config.Anki()", logg.Err(err))
		os.Exit(-1)
	}
	return anki.Config{ExportPrefix: "t2a-", NotesCacheDir: cacheDir}
}
