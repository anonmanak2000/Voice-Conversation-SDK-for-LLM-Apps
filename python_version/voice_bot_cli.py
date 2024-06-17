import argparse
import pyaudio
import time
import threading
import queue
import openai
from deepgram import Deepgram
from typing import Protocol
import asyncio

class InputStream(Protocol):
    def read(self, chunk_size: int):
        pass

class OutputStream(Protocol):
    def write(self, data: bytes):
        pass

class PyAudioInputStream:
    def __init__(self, stream):
        self.stream = stream

    def read(self, chunk_size: int):
        return self.stream.read(chunk_size)

def create_input_stream(sample_rate: int, chunk_size: int):
    p = pyaudio.PyAudio()
    stream = p.open(
        format=pyaudio.paInt16,
        channels=1,
        rate=sample_rate,
        input=True,
        frames_per_buffer=chunk_size,
    )
    return PyAudioInputStream(stream)

class PyAudioOutputStream:
    def __init__(self, stream):
        self.stream = stream

    def write(self, data: bytes):
        self.stream.write(data)

def create_output_stream(sample_rate: int, chunk_size: int):
    p = pyaudio.PyAudio()
    stream = p.open(
        format=pyaudio.paInt16,
        channels=1,
        rate=sample_rate,
        output=True,
        frames_per_buffer=chunk_size,
    )
    return PyAudioOutputStream(stream)

class VoiceConversationSDK:
    def __init__(self):
        self.stt_client = None
        self.tts_client = None
        self.llm_prompt = None
        self.metrics = {
            'stt_time': 0,
            'llm_first_token_time': 0,
            'llm_complete_time': 0,
            'tts_time': 0
        }

    def setup(self, stt_config: dict, tts_config: dict, llm_config: dict):
        # Initialize Deepgram client
        self.stt_client = Deepgram(stt_config['api_key'])
        openai.api_key = llm_config['api_key']

    async def process_speech(self, input_stream: InputStream, output_stream: OutputStream):
        audio_queue = queue.Queue()
        stt_complete_event = threading.Event()
        llm_response_queue = queue.Queue()

        def process_audio():
            while True:
                audio_data = input_stream.read(1024)
                if audio_data:
                    audio_queue.put(audio_data)

        async def stt_process():
            stt_start_time = time.time()
            stt_transcript = []

            async def on_transcript_received(data):
                nonlocal stt_transcript, stt_complete_event
                stt_transcript.append(data['channel']['alternatives'][0]['transcript'])
                if data['is_final']:
                    stt_complete_event.set()

            # Start the live transcription session
            async with self.stt_client.live({'punctuate': True}) as socket:
                socket.on_transcript(on_transcript_received)

                while not stt_complete_event.is_set():
                    if not audio_queue.empty():
                        audio_data = audio_queue.get()
                        await socket.send(audio_data)

                await socket.stop()
            
            stt_end_time = time.time()
            self.metrics['stt_time'] = stt_end_time - stt_start_time
            return ''.join(stt_transcript)

        async def llm_process(transcript, system_prompt):
            llm_start_time = time.time()
            completion = openai.Completion.create(
                engine="davinci",
                prompt=f"{system_prompt}\nUser: {transcript}\nAI:",
                max_tokens=150
            )
            llm_first_token_time = time.time()
            self.metrics['llm_first_token_time'] = llm_first_token_time - llm_start_time
            llm_response = completion.choices[0].text.strip()
            llm_end_time = time.time()
            self.metrics['llm_complete_time'] = llm_end_time - llm_start_time
            return llm_response

        async def tts_process(llm_response):
            tts_start_time = time.time()
            tts_response = openai.Audio.create(
                engine="davinci-codex",
                text=llm_response,
                voice="en_us_male"
            )
            for audio_chunk in tts_response['audio']:
                output_stream.write(audio_chunk)
            tts_end_time = time.time()
            self.metrics['tts_time'] = tts_end_time - tts_start_time

        # Capture system prompt
        print("Please provide the system prompt by speaking into the microphone...")
        audio_thread = threading.Thread(target=process_audio)
        audio_thread.start()
        system_prompt = await stt_process()
        audio_thread.join()

        # Capture user input for conversation
        print("System prompt captured. Now start the main conversation...")
        audio_thread = threading.Thread(target=process_audio)
        audio_thread.start()
        user_transcript = await stt_process()
        audio_thread.join()

        # Process LLM and TTS
        llm_response = await llm_process(user_transcript, system_prompt)
        await tts_process(llm_response)

    def get_performance_metrics(self):
        return self.metrics

def main():
    parser = argparse.ArgumentParser(description="Voice Bot CLI")
    parser.add_argument('--stt-api-key', type=str, required=True, help='API key for Deepgram STT')
    parser.add_argument('--tts-api-key', type=str, required=True, help='API key for OpenAI TTS')
    parser.add_argument('--llm-api-key', type=str, required=True, help='API key for OpenAI LLM')
    
    args = parser.parse_args()
    
    stt_config = {
        'name': 'Deepgram',
        'api_key': args.stt_api_key
    }

    tts_config = {
        'name': 'OpenAI',
        'api_key': args.tts_api_key
    }

    llm_config = {
        'name': 'OpenAI',
        'api_key': args.llm_api_key
    }

    sdk = VoiceConversationSDK()
    sdk.setup(stt_config, tts_config, llm_config)

    input_stream = create_input_stream(sample_rate=16000, chunk_size=1024)
    output_stream = create_output_stream(sample_rate=16000, chunk_size=1024)

    asyncio.run(sdk.process_speech(input_stream, output_stream))
    metrics = sdk.get_performance_metrics()
    print("Performance metrics:", metrics)

if __name__ == "__main__":
    main()
