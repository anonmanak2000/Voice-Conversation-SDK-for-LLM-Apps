package voicebot

import (
	"fmt"
	"log"
	"time"
)

type Config struct {
	Deepgram_API_KEY string
	OPENAI_API_KEY   string
}

type VoiceBotSDK struct {
	config  Config
	metrics metrics
}

type metrics struct {
	STTTime           time.Duration
	LLMFirstTokenTime time.Duration
	LLMCompleteTime   time.Duration
	TTSTime           time.Duration
}

func NewVoiceBotSDK(config Config) *VoiceBotSDK {
	return &VoiceBotSDK{
		config: Config{
			Deepgram_API_KEY: config.Deepgram_API_KEY,
			OPENAI_API_KEY:   config.OPENAI_API_KEY,
		},

		metrics: metrics{},
	}
}

func (sdk *VoiceBotSDK) Process(fileName string) error {
	sttStart := time.Now()

	// Read the WAV file
	data, err := readWavFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read WAV file: %v", err)
	}

	//Start Speech to Text Conversion
	transcript, err := convertSpeechToText(sdk.config.Deepgram_API_KEY, data)
	if err != nil {
		log.Fatalf("Failed to convert speech to text: %v", err)
	}
	log.Printf("Transcript: %s", transcript)
	sttEnd := time.Now()
	sdk.metrics.STTTime = sttEnd.Sub(sttStart)

	llmStart := time.Now()

	//Start LLM query
	response, err := getOpenAIResponse(sdk.config.OPENAI_API_KEY, transcript)
	if err != nil {
		log.Fatalf("Failed to get response from OpenAI: %v", err)
	}
	log.Printf("Response: %s", response)

	llmEnd := time.Now()
	sdk.metrics.LLMFirstTokenTime = llmEnd.Sub(llmStart)
	sdk.metrics.LLMCompleteTime = sdk.metrics.LLMFirstTokenTime

	ttsStart := time.Now()

	//Start Text to Speech Conversion
	outputAudioData, err := convertTextToSpeech(sdk.config.OPENAI_API_KEY, response)
	if err != nil {
		log.Fatalf("Failed to convert text to speech: %v", err)
	}
	// Save the audio data as output.wav
	err = saveWavFile("output.wav", outputAudioData)
	if err != nil {
		log.Fatalf("Failed to save WAV file: %v", err)
	}
	log.Println("output.wav file has been saved successfully")
	ttsEnd := time.Now()
	sdk.metrics.TTSTime = ttsEnd.Sub(ttsStart)

	//Print Metrics for the request
	fmt.Println("Metrics for the library: ")

	currentMetrics := sdk.GetMetrics()
	fmt.Println("Speech to Text Time: ", currentMetrics.STTTime)

	fmt.Println("Time to retreive first token from LLM: ", currentMetrics.LLMFirstTokenTime)

	fmt.Println("Time taken by LLM to complete:  ", currentMetrics.LLMCompleteTime)

	fmt.Println("Text to Speech Time: ", currentMetrics.TTSTime)

	return nil
}

func (sdk *VoiceBotSDK) GetMetrics() metrics {
	return sdk.metrics
}
