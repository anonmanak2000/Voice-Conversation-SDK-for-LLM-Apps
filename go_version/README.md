# Voice Bot CLI Application

## Overview

This is a Command-Line Interface (CLI) application written in Go for a voice bot that processes speech input from a microphone, interacts with a Large Language Model (LLM) for natural language processing, and provides speech output through audio playback.

### Features

- **Speech-to-Text (STT)**: Converts speech input into text using Deepgram's API.
- **Large Language Model (LLM)**: Interacts with OpenAI's GPT-3.5 to generate responses based on user input.
- **Text-to-Speech (TTS)**: Converts text responses from the LLM into speech for audio output.
- **Concurrency**: Uses goroutines to handle microphone input and audio output concurrently.
- **Performance Metrics**: Tracks various metrics such as STT time, LLM response time, and TTS response time.
- **CLI Interface**: Provides a simple command-line interface to start and stop the application.

## Setup

### Prerequisites

- Go (version 1.16+)
- PortAudio library (for audio I/O handling)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/anonmanak2000/Voice-Conversation-SDK-for-LLM-Apps.git
   cd Voice-Conversation-SDK-for-LLM-Apps
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up API keys:

   STT API Key: Obtain from Deepgram for speech-to-text conversion.
   TTS API Key: Obtain from OpenAI for text-to-speech conversion.
   LLM API Key: Obtain from OpenAI for interacting with the Large Language Model.

Store these keys securely and configure them using command-line flags or environment variables.
