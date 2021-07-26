// Package azure provides access to the Azure text to speech REST API
package azure

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/s12chung/text2anki/pkg/synthesizers"
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

	client *http.Client
}

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

// NewAzure returns a new Azure API struct
func NewAzure(apiKey string, region Region) synthesizers.Synthesizer {
	return &Azure{apiKey: apiKey, region: region, client: http.DefaultClient}
}

//nolint:gosec // no creds here
const tokenURL = "https://%v.api.cognitive.microsoft.com/sts/v1.0/issueToken"

// Token returns a API token
func (a *Azure) Token() (string, error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf(tokenURL, a.region), nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	a.addAPIKey(request.Header)

	response, err := a.client.Do(request)
	if err != nil {
		return "", err
	}
	token, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

const textToSpeechURL = "https://%v.tts.speech.microsoft.com/cognitiveservices/v1"
const textToSpeechBody = `
<speak version='1.0' xml:lang='ko-KR'>
<voice xml:lang='ko-KR' xml:gender='Female' name='ko-KR-SunHiNeural'>
%v
</voice>
</speak>
`

// TextToSpeech returns the speech audio of the given string
func (a *Azure) TextToSpeech(s string) ([]byte, error) {
	if err := a.setupToken(); err != nil {
		return nil, err
	}

	body := bytes.NewBufferString(fmt.Sprintf(textToSpeechBody, s))
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf(textToSpeechURL, a.region), body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-type", "application/ssml+xml")
	request.Header.Add("User-Agent", "text2anki")
	request.Header.Add("X-Microsoft-OutputFormat", "audio-24khz-96kbitrate-mono-mp3")
	a.addToken(request.Header)

	response, err := a.client.Do(request)
	if err != nil {
		return nil, err
	}
	speech, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return speech, nil
}

const apiKeyHeader = "Ocp-Apim-Subscription-Key"
const tokenHeader = "Authorization"

func (a *Azure) addAPIKey(header http.Header) {
	header.Add(apiKeyHeader, a.apiKey)
}

func (a *Azure) setupToken() error {
	token, err := a.Token()
	if err != nil {
		return err
	}
	a.token = token
	return nil
}

func (a *Azure) addToken(header http.Header) {
	header.Add(tokenHeader, "Bearer "+a.token)
}
