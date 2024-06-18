package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type DeepgramResponse struct {
	Results struct {
		Channels []struct {
			Alternatives []struct {
				Transcript string `json:"transcript"`
			} `json:"alternatives"`
		} `json:"channels"`
	} `json:"results"`
}

type OpenAIRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	MaxTokens int `json:"max_tokens"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type OpenAITTSRequest struct {
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

type OpenAITTSResponse struct {
	AudioContent string `json:"audio_content"`
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

	fmt.Println("Response from Deepgram API: ", string(body))

	var deepgramResp DeepgramResponse
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
	requestBody, err := json.Marshal(OpenAIRequest{
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

	var openAIResp OpenAIResponse
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

	requestBody, err := json.Marshal(OpenAITTSRequest{
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

	var ttsResponse OpenAITTSResponse
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

func saveWavFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func main() {

	DEEPGRAM_KEY := "70187200ed45f0adca6eca8a3f99ebc7127b1f1d"
	OPENAPI_KEY := "sk-proj-HUnBpcgnhE3DchFIlpZ8T3BlbkFJ7LwVyjN82cWBE93s7OO0"

	os.Setenv("DEEPGRAM_API_KEY", DEEPGRAM_KEY)
	os.Setenv("OPENAI_API_KEY", OPENAPI_KEY)

	// Read the WAV file
	data, err := readWavFile("input.wav")
	if err != nil {
		log.Fatalf("Failed to read WAV file: %v", err)
	}

	// Convert speech to text using Deepgram API
	deepgramAPIKey := os.Getenv("DEEPGRAM_API_KEY")
	transcript, err := convertSpeechToText(deepgramAPIKey, data)
	if err != nil {
		log.Fatalf("Failed to convert speech to text: %v", err)
	}
	log.Printf("Transcript: %s", transcript)

	// Get a meaningful answer from OpenAI API
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	response, err := getOpenAIResponse(openAIAPIKey, transcript)
	if err != nil {
		log.Fatalf("Failed to get response from OpenAI: %v", err)
	}
	log.Printf("Response: %s", response)

	// Convert text to speech
	audioData, err := convertTextToSpeech(openAIAPIKey, response)
	if err != nil {
		log.Fatalf("Failed to convert text to speech: %v", err)
	}

	// Save the audio data as output.wav
	err = saveWavFile("output.wav", audioData)
	if err != nil {
		log.Fatalf("Failed to save WAV file: %v", err)
	}
	log.Println("output.wav has been saved successfully")
}
