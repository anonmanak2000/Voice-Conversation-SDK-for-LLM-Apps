package main

import (
	"anonmanak2000/Voice-Conversation-SDK-for-LLM-Apps/voicebot"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {

	DEEPGRAM_KEY := goDotEnvVariable("DEEPGRAM_KEY")
	OPENAPI_KEY := goDotEnvVariable("OPENAPI_KEY")

	os.Setenv("DEEPGRAM_API_KEY", DEEPGRAM_KEY)
	os.Setenv("OPENAI_API_KEY", OPENAPI_KEY)

	config := voicebot.Config{
		Deepgram_API_KEY: DEEPGRAM_KEY,
		OPENAI_API_KEY:   OPENAPI_KEY,
	}
	voiceBotSDK := voicebot.NewVoiceBotSDK(config)

	voiceBotSDK.Process("input.wav") //Provide filename as an argument

}
