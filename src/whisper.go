package main

import (
	"os"
	"os/exec"
	"path"
)

func Transcribe(content string, dir string) []byte {
	log("Transcribing string")

	trans := exec.Command(
		"whisper",
		"audio.mp3",
		"--model", "tiny.en",
		"--language", "en",
		"--fp16", "False",
		"--device", "cuda",
	)

	trans.Dir = dir

	// the woke left in my fucking software
	trans.Run()

	transcription, _ := os.ReadFile(path.Join(dir, "audio.txt"))

	return transcription
}
