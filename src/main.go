package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
)

func main() {
	u, _ := user.Current()

	runDir := "/home/" + u.Username + "/.sugoi"
	os.MkdirAll(runDir, 0755)
	os.Remove(path.Join(runDir, "audio.mp3"))
	os.Remove(path.Join(runDir, "audio.txt"))

	url := os.Args[1]

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", "audio.mp3",
		url,
	)

	cmd.Dir = runDir

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running yt-dlp: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Downloaded audio")

	trans := exec.Command(
		"whisper",
		"audio.mp3",
		"--model", "tiny.en",
		"--language", "en",
		"--fp16", "False",
		"--device", "cuda",
	)

	trans.Dir = runDir

	// the woke left in my fucking software
	trans.Run()

	content, _ := os.ReadFile(path.Join(runDir, "audio.txt"))

	fmt.Println("Got transcription")

	ollama := "http://127.0.0.1:11434/api/generate"
	payload := map[string]interface{}{
		"prompt": `Summarize the following YouTube video transcript in 3–4 sentences.
Use concise bullet points.
Do not include any introductory or filler phrases (ie: "Here's a summary of the youtube video:").
Do not reference the fact that it's a video or a transcript.
Focus only on the core events or ideas, like a newsreader delivering headlines.
Do not use markdown or any funky indentation.
\n` + string(content),
		"model":  "gemma3:4b",
		"stream": false,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(ollama, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	type Response struct {
		Content string `json:"response"`
	}

	parse := Response{}

	json.Unmarshal(responseBody, &parse)

	fmt.Print(parse.Content)
}
