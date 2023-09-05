// Package azure provides access to the Azure text to speech REST API
package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/s12chung/text2anki/pkg/synthesizer"
)

// GetAPIKeyFromEnv gets the API key from the default ENV var
func GetAPIKeyFromEnv() string {
	return os.Getenv("AZURE_SPEECH_API_KEY")
}

// Azure is a wrapper of the Azure Text to Speech API
type Azure struct {
	apiKey string
	region Region

	token string

	client       *http.Client
	cache        map[string][]byte
	requestCount uint
}

// requestLimit is the limit of requests tested from the azure API
const requestLimit = 15

// Region are region identifiers for the API
type Region string

// Regions, see: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech
const (
	CentralUSRegion        Region = "centralus"
	EastUSRegion           Region = "eastus"
	EastUS2Region          Region = "eastus2"
	NorthCentralUS2Region  Region = "northcentralus"
	SouthCentralUSRegion   Region = "southcentralus"
	WestCentralUSRegion    Region = "westcentralus"
	WestUSRegion           Region = "westus"
	WestUS2Region          Region = "westus2"
	CanadaCentralRegion    Region = "canadacentral"
	BrazilSouthRegion      Region = "brazilsouth"
	EastAsiaRegion         Region = "eastasia"
	SouthEastAsiaRegion    Region = "southeastasia"
	AustraliaEastRegion    Region = "australiaeast"
	CentralIndiaRegion     Region = "centralindia"
	JapanEastRegion        Region = "japaneast"
	JapanWestRegion        Region = "japanwest"
	KoreaCentralRegion     Region = "koreacentral"
	NorthEuropeRegion      Region = "northeurope"
	WestEuropeRegion       Region = "westeurope"
	FranceCentralRegion    Region = "francecentral"
	SwitzerlandNorthRegion Region = "switzerlandnorth"
	UKSouthRegion          Region = "uksouth"
)

// New returns a new Azure API struct
func New(apiKey string, region Region) synthesizer.Synthesizer {
	return &Azure{apiKey: apiKey, region: region, client: http.DefaultClient, cache: map[string][]byte{}}
}

// SourceName is the source name of this synthesizer
const SourceName = "Azure Text-to-speech REST API"

// SourceName returns the source name of this synthesizer
func (a *Azure) SourceName() string {
	return SourceName
}

//nolint:gosec // no creds here
const tokenURL = "https://%v.api.cognitive.microsoft.com/sts/v1.0/issueToken"

// Token returns a API token
func (a *Azure) Token(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(tokenURL, a.region), nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	a.addAPIKey(request.Header)

	resp, err := a.client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //nolint:errcheck // failing is ok
	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

const textToSpeechURL = "https://%v.tts.speech.microsoft.com/cognitiveservices/v1"
const textToSpeechBody = `
<speak version='1.0' xml:lang='ko-KR'>
<voice xml:lang='ko-KR' xml:gender='Female' name='ko-KR-SunHiNeural'>
<prosody rate="0.75">
%v
</prosody>
</voice>
</speak>
`

// TextToSpeech returns the speech audio of the given string
func (a *Azure) TextToSpeech(ctx context.Context, s string) ([]byte, error) {
	reqBodyString := fmt.Sprintf(textToSpeechBody, s)
	if bytes, exists := a.cache[reqBodyString]; exists {
		return bytes, nil
	}

	if err := a.setupToken(ctx); err != nil {
		return nil, err
	}

	reqBody := bytes.NewBufferString(reqBodyString)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(textToSpeechURL, a.region), reqBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-type", "application/ssml+xml")
	request.Header.Add("User-Agent", "text2anki")
	request.Header.Add("X-Microsoft-OutputFormat", "audio-24khz-96kbitrate-mono-mp3")
	a.addToken(request.Header)

	if a.requestCount != 0 && a.requestCount%requestLimit == 0 {
		time.Sleep(60 * time.Second)
	}
	a.requestCount++

	resp, err := a.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck // failing is ok
	if resp.StatusCode != http.StatusOK {
		var body []byte
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			body = nil
		}
		return nil, fmt.Errorf("returns a non-200 status code: %v (%v) with body: %v",
			resp.StatusCode, resp.Status, string(body))
	}
	speech, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	a.cache[reqBodyString] = speech
	return speech, nil
}

//nolint:gosec // not a credential
const apiKeyHeader = "Ocp-Apim-Subscription-Key"
const tokenHeader = "Authorization"

func (a *Azure) addAPIKey(header http.Header) {
	header.Add(apiKeyHeader, a.apiKey)
}

func (a *Azure) setupToken(ctx context.Context) error {
	token, err := a.Token(ctx)
	if err != nil {
		return err
	}
	a.token = token
	return nil
}

func (a *Azure) addToken(header http.Header) {
	header.Add(tokenHeader, "Bearer "+a.token)
}

// SetClient sets the client for API requests
func (a *Azure) SetClient(c *http.Client) {
	a.client = c
}
