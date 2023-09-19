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
func GetAPIKeyFromEnv() string { return os.Getenv("AZURE_SPEECH_API_KEY") }

// Azure is a wrapper of the Azure Text to Speech API
type Azure struct {
	apiKey string
	region Region

	client       *http.Client
	cache        map[string][]byte
	requestCount uint
}

// requestLimit is the limit of requests tested from the azure API per minute
// https://learn.microsoft.com/en-us/azure/ai-services/speech-service/speech-services-quotas-and-limits#real-time-text-to-speech
const requestLimit = 20

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
func (a *Azure) SourceName() string { return SourceName }

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
func (a *Azure) TextToSpeech(ctx context.Context, text string) ([]byte, error) {
	reqBodyString := fmt.Sprintf(textToSpeechBody, text)
	if bytes, exists := a.cache[reqBodyString]; exists {
		return bytes, nil
	}

	reqBody := bytes.NewBufferString(reqBodyString)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(textToSpeechURL, a.region), reqBody)
	if err != nil {
		return nil, err
	}
	if err := a.addToken(ctx, request.Header); err != nil {
		return nil, err
	}
	request.Header.Add("Content-type", "application/ssml+xml")
	request.Header.Add("User-Agent", "text2anki")
	request.Header.Add("X-Microsoft-OutputFormat", "audio-24khz-96kbitrate-mono-mp3")

	if a.requestCount != 0 && a.requestCount%requestLimit == 0 {
		time.Sleep(65 * time.Second) // 5 second padding
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

func (a *Azure) addAPIKey(header http.Header) { header.Add(apiKeyHeader, a.apiKey) }

func (a *Azure) addToken(ctx context.Context, header http.Header) error {
	token, err := a.Token(ctx)
	if err != nil {
		return err
	}
	header.Add(tokenHeader, "Bearer "+token)
	return nil
}

// SetClient sets the client for API requests
func (a *Azure) SetClient(c *http.Client) { a.client = c }
