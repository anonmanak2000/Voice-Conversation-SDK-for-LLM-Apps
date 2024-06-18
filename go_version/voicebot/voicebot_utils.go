package voicebot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type deepgramResponse struct {
	Results struct {
		Channels []struct {
			Alternatives []struct {
				Transcript string `json:"transcript"`
			} `json:"alternatives"`
		} `json:"channels"`
	} `json:"results"`
}

type openAIRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	MaxTokens int `json:"max_tokens"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openAITTSRequest struct {
	Input struct {
		Text string `json:"text"`
	} `json:"input"`
	Voice struct {
		LanguageCode string `json:"language_code"`
		Name         string `json:"name"`
	} `json:"voice"`
	AudioConfig struct {
		AudioEncoding string `json:"audio_encoding"`
	} `json:"audio_config"`
}

type openAITTSResponse struct {
	AudioContent string `json:"audio_content"`
}

func saveWavFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func readWavFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func convertSpeechToText(apiKey string, audioData []byte) (string, error) {
	url := "https://api.deepgram.com/v1/listen"
	req, err := http.NewRequest("POST", url, bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Token "+apiKey)
	req.Header.Set("Content-Type", "audio/wav")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var deepgramResp deepgramResponse
	err = json.Unmarshal(body, &deepgramResp)
	if err != nil {
		return "", err
	}

	if len(deepgramResp.Results.Channels) > 0 && len(deepgramResp.Results.Channels[0].Alternatives) > 0 {
		return deepgramResp.Results.Channels[0].Alternatives[0].Transcript, nil
	}

	return "", nil
}

func getOpenAIResponse(apiKey, prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	requestBody, err := json.Marshal(openAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 150,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Response from OpenAI: ", string(body))

	var openAIResp openAIResponse
	err = json.Unmarshal(body, &openAIResp)
	if err != nil {
		return "", err
	}

	if len(openAIResp.Choices) > 0 {
		return strings.TrimSpace(openAIResp.Choices[0].Message.Content), nil
	}

	return "", nil
}
func convertTextToSpeech(apiKey, text string) ([]byte, error) {
	url := "https://api.openai.com/v1/text-to-speech"

	requestBody, err := json.Marshal(openAITTSRequest{
		Input: struct {
			Text string `json:"text"`
		}{Text: text},
		Voice: struct {
			LanguageCode string `json:"language_code"`
			Name         string `json:"name"`
		}{LanguageCode: "en-US", Name: "en-US-Wavenet-D"},
		AudioConfig: struct {
			AudioEncoding string `json:"audio_encoding"`
		}{AudioEncoding: "LINEAR16"},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ttsResponse openAITTSResponse
	err = json.Unmarshal(body, &ttsResponse)
	if err != nil {
		return nil, err
	}

	audioData, err := base64.StdEncoding.DecodeString(ttsResponse.AudioContent)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}
