package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GenerateSummary(content []byte) string {
	ollama := "http://127.0.0.1:11434/api/generate"
	payload := map[string]interface{}{
		"prompt": `Summarize the following content in 3–4 sentences.
	Use concise bullet points.
	Do not include any introductory or filler phrases.
	Do not reference the source or format of the content.
	Focus only on the core events or ideas, like a newsreader delivering headlines.
	Do not use markdown or any non standard indentation.
	\n` + string(content),
		"model":  "gemma3:4b",
		"stream": false,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(ollama, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	type Response struct {
		Content string `json:"response"`
	}

	parse := Response{}

	json.Unmarshal(responseBody, &parse)

	return parse.Content
}
