package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type SpeechToTextWorker struct {
	option *Option
}

func NewSpeechToTextWorker(option *Option) *SpeechToTextWorker {
	return &SpeechToTextWorker{
		option: option,
	}
}

func (w *SpeechToTextWorker) HasMagicalSpel(file *os.File) (bool, error) {

	token, err := w.issueToken()
	if err != nil {
		return false, err
	}

	apiURL := w.getRecognizeURL()

	b := bufio.NewReader(file)
	req, err := http.NewRequest("POST", apiURL, b)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "audio/wav; samplerate=8192")

	client := &http.Client{}
	client.Timeout = 30 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		result := RecognitionResult{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return false, err
		}

		return result.HasMagicSpel(w.option.MagicalSpel), nil
	}

	if resp.StatusCode == 400 || resp.StatusCode == 415 || resp.StatusCode == 500 {
		body, _ := ioutil.ReadAll(resp.Body)

		return false, fmt.Errorf("Failed on calling recognition API: %s", body)
	}

	return false, fmt.Errorf("Unknown Error Occurred: recognizing, Check the key , Status : " + resp.Status)

}

func (w *SpeechToTextWorker) issueToken() (string, error) {
	// REST APIでazureトークン取得
	apiURL := "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"

	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", w.option.AzureSubscriptionKey)
	req.Header.Set("Content-Length", "0")

	client := &http.Client{}
	client.Timeout = 30 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), nil
	}
	if resp.StatusCode == 400 || resp.StatusCode == 415 || resp.StatusCode == 500 {
		body, _ := ioutil.ReadAll(resp.Body)

		return "", fmt.Errorf("Failed on getting token : %s", body)
	}

	return "", fmt.Errorf("Unknown Error Occurred: issueToken , Check the key , Status : " + resp.Status)
}

func (w *SpeechToTextWorker) getRecognizeURL() string {
	urlBase := "https://speech.platform.bing.com/recognize"
	version := "3.0"
	requestid := uuid.NewV4()
	appID := "D4D52672-91D7-4C74-8AD8-42B1D98141A5"
	format := "json"
	locale := "ja-JP"
	incetanceid := uuid.NewV4()
	device := "sacloud-balus"

	return fmt.Sprintf("%s?version=%s&requestid=%s&appID=%s&format=%s&locale=%s&device.os=%s&scenarios=ulm&instanceid=%s",
		urlBase, version, requestid, appID, format, locale, device, incetanceid)
}

// ----------------------------------------------------------------------------

type RecognitionResult struct {
	Header struct {
		Lexical    string `json:"lexical"`
		Name       string `json:"name"`
		Properties struct {
			Requestid string `json:"requestid"`
		} `json:"properties"`
		Scenario string `json:"scenario"`
		Status   string `json:"status"`
	} `json:"header"`
	Results []struct {
		Lexical    string `json:"lexical"`
		Name       string `json:"name"`
		Properties struct {
			Highconf string `json:"HIGHCONF"`
		} `json:"properties"`
		Tokens []struct {
			Lexical       string `json:"lexical"`
			Pronunciation string `json:"pronunciation"`
			Token         string `json:"token"`
		} `json:"tokens"`
	} `json:"results"`
	Version string `json:"version"`
}

func (r *RecognitionResult) HasMagicSpel(spel string) bool {
	if r.Header.Status != "success" {
		return false
	}

	for _, res := range r.Results {
		if res.Lexical == spel {
			return true
		}
		for _, token := range res.Tokens {
			if token.Token == spel || token.Lexical == spel {
				return true
			}
		}
	}
	return false
}
