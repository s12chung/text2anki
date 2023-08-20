package api

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/dictionary/krdict"
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
)

// Config contains config settings for the API
type Config struct {
	TokenizerType
	DictionaryType
	StorageConfig StorageConfig
}

// NewRoutes is the routes used by the API
func NewRoutes(config Config) Routes {
	routes := Routes{
		Dictionary:  Dictionary(config.DictionaryType),
		Synthesizer: Synthesizer(),
		TextTokenizer: db.TextTokenizer{
			Parser:       Parser(),
			Tokenizer:    Tokenizer(config.TokenizerType),
			CleanSpeaker: true,
		},
		Storage: StorageFromConfig(config.StorageConfig),
	}
	db.SetDBStorage(routes.Storage.DBStorage)
	return routes
}

// Parser returns the default Parser
func Parser() text.Parser {
	return text.NewParser(text.Korean, text.English)
}

// Synthesizer returns the default Synthesizer
func Synthesizer() synthesizer.Synthesizer {
	return azure.New(azure.GetAPIKeyFromEnv(), azure.EastUSRegion)
}

// TokenizerType is an enum of tokenizer types
type TokenizerType int

const (
	// TokenizerKhaiii picks the Khaiii tokenizer
	TokenizerKhaiii TokenizerType = iota
	// TokenizerKomoran picks the Komoran tokenizer
	TokenizerKomoran
)

// Tokenizer returns the default Tokenizer
func Tokenizer(tokenizerType TokenizerType) tokenizer.Tokenizer {
	switch tokenizerType {
	case TokenizerKomoran:
		return komoran.New()
	case TokenizerKhaiii:
		fallthrough
	default:
		return khaiii.New()
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
		if err := db.SetDB("db/data.sqlite3"); err != nil {
			fmt.Println("failure to SetDB()\n", err)
			os.Exit(-1)
		}
		return krdict.New(db.DB())
	}
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

// StorageFromConfig returns a storage from the given config
func StorageFromConfig(config StorageConfig) Storage {
	var api storage.API
	var storer storage.Storer
	var err error
	switch config.StorageType {
	case StorageLocalStore:
		fallthrough
	default:
		var ls localstore.API
		ls, err = LocalStoreAPI(config.LocalStoreConfig)
		api = ls
		storer = ls
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return Storage{DBStorage: storage.NewDBStorage(api, config.UUIDGenerator), Storer: storer}
}

// LocalStoreConfig defines the config for localstore
type LocalStoreConfig struct {
	Origin        string
	KeyBasePath   string
	EncryptorPath string
}

var localStoreConfigValidator = firm.NewStructValidator(firm.RuleMap{
	"Origin":        {rule.Presence{}},
	"KeyBasePath":   {rule.Presence{}},
	"EncryptorPath": {rule.Presence{}},
})

const localstoreKey = "localstore.key"
const storageURLPath = "/storage"

// LocalStoreAPI returns a localstore.API
func LocalStoreAPI(config LocalStoreConfig) (localstore.API, error) {
	if !strings.HasSuffix(config.Origin, "/") {
		config.Origin += "/"
	}
	config.Origin += storageURLPath[1:]

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
