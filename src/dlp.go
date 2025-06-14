package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func DownloadAudio(url string, dir string) string {
	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", "audio.mp3",
		"--postprocessor-args", "-ar 16000",
		"--no-playlist",
		url,
	)

	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running yt-dlp: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Downloaded audio")

	return path.Join(dir, "audio.mp4")
}

func IsSupported(url string) bool {
	cmd := exec.Command("yt-dlp", "--dump-json", "--simulate", "--skip-download", url, "--no-playlist", "--no-warnings", "--no-call-home")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false
	}

	return strings.Contains(string(output), `"url"`)
}
