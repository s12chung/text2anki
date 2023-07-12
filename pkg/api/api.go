// Package api contains the routes for the api
package api

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/dictionary/krdict"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/synthesizers"
	"github.com/s12chung/text2anki/pkg/synthesizers/azure"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/khaiii"
	"github.com/s12chung/text2anki/pkg/tokenizers/komoran"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

// Config contains config settings for the API
type Config struct {
	TokenizerType
	DictionaryType
	SignerConfig
}

// NewRoutes is the routes used by the API
func NewRoutes(config Config) Routes {
	return Routes{
		Dictionary:  Dictionary(config.DictionaryType),
		Synthesizer: Synthesizer(),
		TextTokenizer: db.TextTokenizer{
			Parser:       Parser(),
			Tokenizer:    Tokenizer(config.TokenizerType),
			CleanSpeaker: true,
		},
		Signer: Signer(config.SignerConfig),
	}
}

// Parser returns the default Parser
func Parser() text.Parser {
	return text.NewParser(text.Korean, text.English)
}

// Synthesizer returns the default Synthesizer
func Synthesizer() synthesizers.Synthesizer {
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
func Tokenizer(tokenizerType TokenizerType) tokenizers.Tokenizer {
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

// SignerType defines what signer type to for file storage
type SignerType int

const (
	// SignerLocalStore picks the local file store
	SignerLocalStore SignerType = iota
)

// SignerConfig configures the storage signer
type SignerConfig struct {
	SignerType
	LocalStoreConfig LocalStoreConfig
}

// Signer returns a storage signer
func Signer(config SignerConfig) storage.Signer {
	var api storage.API
	var err error
	switch config.SignerType {
	case SignerLocalStore:
		fallthrough
	default:
		api, err = LocalStoreAPI(config.LocalStoreConfig)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return storage.NewSigner(api)
}

// LocalStoreConfig defines the config for localstore
type LocalStoreConfig struct {
	Origin   string
	BaseBath string
	KeyPath  string
}

var localStoreConfigRuleMap = firm.RuleMap{
	"Origin":   {rule.Presence{}},
	"BaseBath": {rule.Presence{}},
	"KeyPath":  {rule.Presence{}},
}

const localstoreKey = "localstore.key"

// LocalStoreAPI returns a localstore.API
func LocalStoreAPI(config LocalStoreConfig) (localstore.API, error) {
	// LocalStoreAPI is called when declaring package level vars (before init()), this ensures the definition works
	result := firm.NewStructValidator(localStoreConfigRuleMap).Validate(config)
	if !result.IsValid() {
		return localstore.API{}, fmt.Errorf(result.ErrorMap().String())
	}
	encryptor, err := localstore.NewAESEncryptorFromFile(path.Join(config.KeyPath, localstoreKey))
	if err != nil {
		return localstore.API{}, err
	}
	return localstore.NewAPI(config.Origin, config.BaseBath, encryptor), nil
}

// Routes contains the routes used for the api
type Routes struct {
	Dictionary    dictionary.Dictionary
	Synthesizer   synthesizers.Synthesizer
	TextTokenizer db.TextTokenizer
	Signer        storage.Signer
}

// Setup sets up the routes
func (rs Routes) Setup() error {
	return rs.TextTokenizer.Setup()
}

// Cleanup cleans up the routes
func (rs Routes) Cleanup() error {
	return rs.TextTokenizer.Cleanup()
}

// Router returns the router with all the routes set
func (rs Routes) Router() chi.Router {
	r := chi.NewRouter()
	r.Route("/sources", func(r chi.Router) {
		r.Get("/", httptyped.RespondTypedJSONWrap(rs.SourceIndex))
		r.Post("/", httptyped.RespondTypedJSONWrap(rs.SourceCreate))
		r.Get("/sign_parts", httptyped.RespondTypedJSONWrap(rs.SignParts))

		r.Route("/{sourceID}", func(r chi.Router) {
			r.Use(httputil.RequestWrap(SourceCtx))
			r.Get("/", httptyped.RespondTypedJSONWrap(rs.SourceGet))
			r.Patch("/", httptyped.RespondTypedJSONWrap(rs.SourceUpdate))
			r.Delete("/", httptyped.RespondTypedJSONWrap(rs.SourceDestroy))
		})
	})
	r.Route("/terms", func(r chi.Router) {
		r.Get("/search", httptyped.RespondTypedJSONWrap(rs.TermsSearch))
	})
	r.Route("/notes", func(r chi.Router) {
		r.Post("/", httptyped.RespondTypedJSONWrap(rs.NoteCreate))
	})
	return r
}

func extractAndValidate(r *http.Request, req any) *httputil.HTTPError {
	if httpError := httputil.ExtractJSON(r, req); httpError != nil {
		return httpError
	}
	result := firm.Validate(req)
	if !result.IsValid() {
		return httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf(result.ErrorMap().String()))
	}
	return nil
}
