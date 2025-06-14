package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"

	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
	"golang.org/x/term"
)

type args struct {
	Positional string `arg:"positional" help:"The main data (string, url, file path)"`
	Quiet      bool   `arg:"-q" help:"Don't echo anything except the summary"`
}

var log func(content string)

func main() {
	u, _ := user.Current()

	runDir := "/home/" + u.Username + "/.sugoi"
	os.MkdirAll(runDir, 0755)
	os.Remove(path.Join(runDir, "audio.mp3"))
	os.Remove(path.Join(runDir, "audio.txt"))

	hasPipe := !term.IsTerminal(int(os.Stdin.Fd()))

	var a args

	arg.MustParse(&a)

	log = func(content string) {
		if !a.Quiet {
			fmt.Println(content)
		}
	}

	var pipe []byte

	if hasPipe {
		pipe, _ = io.ReadAll(os.Stdin)
	}

	log("Checking args....")

	result := ""

	if hasPipe {
		log("Summarising pipe")
		result = GenerateSummary(pipe)
	} else if IsURL(a.Positional) {
		if IsSupported(a.Positional) {
			log("Transcribing " + color.YellowString(a.Positional) + " with yt-dlp")
			// TODO: implement yt-dlp calling here

			DownloadAudio(a.Positional, runDir)
			trans := Transcribe(a.Positional, runDir)
			result = GenerateSummary(trans)
		} else {
			log("[Error] URL was passed, but it's not supported by yt-dlp!")
			return
		}
	} else {
		info, stat := os.Stat(a.Positional)

		// assume it's a string
		if os.IsNotExist(stat) {
			log("Summarising string " + color.YellowString(`"`+a.Positional+`"`))

			result = GenerateSummary([]byte(a.Positional))
		} else if !info.IsDir() {
			// it's a file
			log("Reading file " + color.YellowString(a.Positional) + " and transcribing")

			data, err := os.ReadFile(a.Positional)

			if err != nil {
				log("Failed reading supplied file path contents")
				return
			}

			result = GenerateSummary(data)
		} else {
			log(color.RedString("Directory paths are not supported!"))
			return
		}
	}

	fmt.Println("\n" + color.GreenString(result) + "\n")
}
