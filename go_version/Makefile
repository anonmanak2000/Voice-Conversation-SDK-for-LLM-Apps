DEEPGRAM_KEY=
OPENAI_KEY=

build:
	go build -o voicebot-cli main.go

run:build
	./voicebot-cli

cli:build
	./voicebot-cli DEEPGRAM_KEY=${DEEPGRAM_KEY} OPENAI_KEY=${OPENAI_KEY}