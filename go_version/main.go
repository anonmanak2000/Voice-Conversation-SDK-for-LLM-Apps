package main

import (
	"anonmanak2000/Voice-Conversation-SDK-for-LLM-Apps/voicebot"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
		return ""
	}

	return os.Getenv(key)
}

func main() {

	DEEPGRAM_KEY := goDotEnvVariable("DEEPGRAM_KEY")
	OPENAPI_KEY := goDotEnvVariable("OPENAPI_KEY")

	if DEEPGRAM_KEY == "" && OPENAPI_KEY == "" {
		DEEPGRAM_KEY = *flag.String("DEEPGRAM_KEY", "", "Provide Deepgram API Key")
		OPENAPI_KEY = *flag.String("OPENAPI_KEY", "", "Provide Open API Key")
		flag.Parse()
	}

	config := voicebot.Config{
		Deepgram_API_KEY: DEEPGRAM_KEY,
		OPENAI_API_KEY:   OPENAPI_KEY,
	}
	voiceBotSDK := voicebot.NewVoiceBotSDK(config)

	voiceBotSDK.Process("input.wav") //Provide filename as an argument

}
