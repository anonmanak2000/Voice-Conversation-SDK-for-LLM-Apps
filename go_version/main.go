// package main

// import (
// 	"anonmanak2000/Voice-Conversation-SDK-for-LLM-Apps/voicebot"
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/go-audio/audio"
// 	"github.com/go-audio/wav"
// )

// var (
// 	sttAPIKey string = "" // Replace with your API keys
// 	ttsAPIKey string = "" // Replace with your API keys
// 	llmAPIKey string = "" // Replace with your API keys
// )

// func main() {
// 	// Initialize VoiceBotSDK
// 	config := voicebot.Config{
// 		STTApiKey: sttAPIKey,
// 		TTSApiKey: ttsAPIKey,
// 		LLMApiKey: llmAPIKey,
// 	}
// 	sdk := voicebot.NewVoiceBotSDK(config)

// 	// Setup interrupt handler to gracefully stop
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
// 	defer signal.Stop(stop)

// 	// Read input .wav file
// 	wavFile := "input.wav"
// 	audioData, sampleRate, err := readWavFile(wavFile)
// 	if err != nil {
// 		log.Fatalf("Error reading .wav file: %v", err)
// 	}

// 	// Process audio input
// 	output := &bytes.Buffer{}
// 	err = sdk.ProcessSpeech(audioData, output)
// 	if err != nil {
// 		log.Fatalf("Error processing speech: %v", err)
// 	}

// 	// Write output .wav file
// 	outFile := "output.wav"
// 	err = writeWavFile(outFile, audioData, sampleRate)
// 	if err != nil {
// 		log.Fatalf("Error writing output .wav file: %v", err)
// 	}

// 	fmt.Printf("Output saved to %s\n", outFile)
// }

// func readWavFile(filename string) ([]int16, int, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	defer file.Close()

// 	decoder := wav.NewDecoder(file)
// 	if !decoder.IsValidFile() {
// 		return nil, 0, fmt.Errorf("invalid WAV file")
// 	}

// 	buf, err := decoder.FullPCMBuffer()
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	audioData := make([]int16, buf.NumFrames()*buf.Format.NumChannels)
// 	idx := 0
// 	for i := 0; i < buf.NumFrames(); i++ {
// 		for ch := 0; ch < buf.Format.NumChannels; ch++ {
// 			audioData[idx] = int16(binary.LittleEndian.Uint16([]byte{
// 				byte(buf.Data[idx*2]),
// 				byte(buf.Data[idx*2+1]),
// 			}))
// 			idx++
// 		}
// 	}

// 	return audioData, int(buf.Format.SampleRate), nil
// }

// func writeWavFile(filename string, audioData []int16, sampleRate int) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	encoder := wav.NewEncoder(file, sampleRate, 16, 1, 1)
// 	if err := encoder.Write(&audio.IntBuffer{
// 		Data:   BytesToIntSlice(intBuffer(audioData), 2),
// 		Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1},
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func intBuffer(data []int16) []byte {
// 	buf := make([]byte, len(data)*2)
// 	for i := range data {
// 		binary.LittleEndian.PutUint16(buf[i*2:], uint16(data[i]))
// 	}
// 	return buf
// }

// func BytesToIntSlice(bytes []byte, size int) []int {

// 	intSize := int(size)
// 	intData := make([]int, len(bytes)/intSize)

// 	for i := 0; i < len(bytes); i += intSize {

// 		intData[i/intSize] = int(int16(binary.LittleEndian.Uint16(bytes[i : i+intSize])))

// 	}

// 	return intData
// }
