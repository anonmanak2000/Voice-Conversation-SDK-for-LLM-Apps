# Voice Bot CLI Application

## Overview

This is a simple application written in Go for a voice bot that processes speech input from an audio file (.wav), interacts with a Large Language Model (LLM) for natural language processing, and provides an output audio file.

### Features

- **Speech-to-Text (STT)**: Converts speech input into text using Deepgram's API.
- **Large Language Model (LLM)**: Interacts with OpenAI's GPT-3.5 to generate responses based on user input.
- **Text-to-Speech (TTS)**: Converts text responses from the LLM into speech for audio output.
- **Performance Metrics**: Tracks various metrics such as STT time, LLM response time, and TTS response time.

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

3. There are 2 options to setup API Keys:

   a. Create an .env file to store the API Keys
   
   b. Pass API keys as command-line arguments. Update Makefile to add your API Keys

   Deepgram API Key: Obtain from Deepgram for speech-to-text conversion.
   OpenAI API Key: Obtain from OpenAI for text-to-speech conversion.
   OpenAI API Key: Obtain from OpenAI for interacting with the Large Language Model.

Store these keys securely and configure them using command-line flags or environment variables.

4. Usage

   a. If you have created .env file to set up API keys, use following

   ```bash
   make run
   ```
   
   b. If you want to use API Keys via command-line

   ```bash
   make cli
   ```

### Performance Metrics

The application tracks the following performance metrics:

- STT Time: Duration for speech-to-text conversion.
- LLM Response Time: Duration from input reception to LLM first token response.
- TTS Time: Duration for text-to-speech conversion.
