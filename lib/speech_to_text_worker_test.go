package lib

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestSpeechToTextWorker_HasMagicalSpel(t *testing.T) {

	option := NewOption()
	option.AzureSubscriptionKey = os.Getenv("AZURE_SUBSCRIPTION_KEY")
	if option.AzureSubscriptionKey == "" {
		log.Println("Please Set ENV 'AZURE_SUBSCRIPTION_KEY'")
		os.Exit(0) // exit normal
	}

	file, err := os.Open("test.wav")
	if !assert.NoError(t, err) {
		return
	}

	worker := NewSpeechToTextWorker(option)

	res, err := worker.HasMagicalSpel(file)
	assert.NoError(t, err)
	assert.True(t, res)

}
