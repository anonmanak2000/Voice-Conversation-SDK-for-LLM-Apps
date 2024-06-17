package voicebot

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Config struct {
	STTApiKey string
	TTSApiKey string
	LLMApiKey string
}

type VoiceBotSDK struct {
	sttApiKey string
	ttsApiKey string
	llmApiKey string
	metrics   Metrics
}

type Metrics struct {
	STTTime           time.Duration
	LLMFirstTokenTime time.Duration
	LLMCompleteTime   time.Duration
	TTSTime           time.Duration
}

func NewVoiceBotSDK(config Config) *VoiceBotSDK {
	return &VoiceBotSDK{
		sttApiKey: config.STTApiKey,
		ttsApiKey: config.TTSApiKey,
		llmApiKey: config.LLMApiKey,
		metrics:   Metrics{},
	}
}

func (sdk *VoiceBotSDK) ProcessSpeech(audioData []int16, output io.Writer) error {
	sttStart := time.Now()
	transcript, err := sdk.sttToText(audioData)
	if err != nil {
		return err
	}
	sttEnd := time.Now()
	sdk.metrics.STTTime = sttEnd.Sub(sttStart)

	llmStart := time.Now()
	prompt := "System prompt captured. Now start the main conversation: " + transcript
	llmResponse, err := sdk.queryLLM(prompt)
	if err != nil {
		return err
	}
	llmEnd := time.Now()
	sdk.metrics.LLMFirstTokenTime = llmEnd.Sub(llmStart)
	sdk.metrics.LLMCompleteTime = sdk.metrics.LLMFirstTokenTime

	ttsStart := time.Now()
	ttsResponse, err := sdk.textToSpeech(llmResponse)
	if err != nil {
		return err
	}

	// Play audio output through the provided output writer
	_, err = output.Write(ttsResponse)
	if err != nil {
		return err
	}
	ttsEnd := time.Now()
	sdk.metrics.TTSTime = ttsEnd.Sub(ttsStart)

	return nil
}

func (sdk *VoiceBotSDK) sttToText(audioData []int16) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"data": audioData,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.deepgram.com/v1/listen", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+sdk.sttApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	transcript := result["channels"].(map[string]interface{})["alternatives"].([]interface{})[0].(map[string]interface{})["text"].(string)
	return transcript, nil
}

func (sdk *VoiceBotSDK) queryLLM(prompt string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"model":      "text-davinci-003",
		"prompt":     prompt,
		"max_tokens": "150",
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sdk.llmApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	llmResponse := result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
	return llmResponse, nil
}

func (sdk *VoiceBotSDK) textToSpeech(text string) ([]byte, error) {
	requestBody, err := json.Marshal(map[string]string{
		"text":  text,
		"voice": "en_us_male",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/voices", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sdk.ttsApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}

func (sdk *VoiceBotSDK) GetMetrics() Metrics {
	return sdk.metrics
}
