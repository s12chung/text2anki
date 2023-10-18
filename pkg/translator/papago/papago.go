// Package papago provides access to the Papago translation REST API
package papago

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/s12chung/text2anki/pkg/util/httputil"
)

// GetClientCredentialsFromEnv gets the credentials from the default ENV var
func GetClientCredentialsFromEnv() (string, string) {
	return os.Getenv("PAPAGO_CLIENT_ID"), os.Getenv("PAPAGO_CLIENT_SECRET")
}

// Papago provides access to the Papago translation REST API
type Papago struct {
	clientID     string
	clientSecret string

	client *http.Client
	cache  map[string]string
}

// New returns a new Papago
func New(clientID, clientSecret string) *Papago {
	return &Papago{clientID: clientID, clientSecret: clientSecret, client: http.DefaultClient, cache: map[string]string{}}
}

const apiURL = "https://naveropenapi.apigw.ntruss.com/nmt/v1/translation"
const clientIDHeader = "X-Ncp-Apigw-Api-Key-Id"  //nolint:gosec // not a credential
const clientSecretHeader = "X-Ncp-Apigw-Api-Key" //nolint:gosec // not a credential

type translateRequest struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Text   string `json:"text"`
}

type translateResponse struct {
	Message struct {
		Result struct {
			SrcLangType    string `json:"srcLangType"`    //nolint:tagliatelle // API response
			TarLangType    string `json:"tarLangType"`    //nolint:tagliatelle // API response
			TranslatedText string `json:"translatedText"` //nolint:tagliatelle // API response
		} `json:"result"`
	} `json:"message"`
}

// Translate translates the string
func (p *Papago) Translate(ctx context.Context, s string) (string, error) {
	s = strings.TrimSpace(s)
	if translation, exists := p.cache[s]; exists {
		return translation, nil
	}

	reqBody, err := json.Marshal(translateRequest{Source: "ko", Target: "en", Text: s})
	if err != nil {
		return "", err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-type", "application/json")
	request.Header.Add(clientIDHeader, p.clientID)
	request.Header.Add(clientSecretHeader, p.clientSecret)

	respBytes, err := httputil.DoFor200(p.client, request)
	if err != nil {
		return "", err
	}
	resp := translateResponse{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", err
	}
	translation := resp.Message.Result.TranslatedText
	p.cache[s] = translation
	return translation, nil
}

// SetClient sets the client for API requests
func (p *Papago) SetClient(c *http.Client) { p.client = c }
