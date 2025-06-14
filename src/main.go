package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"

	"github.com/fatih/color"
	"golang.org/x/term"
)

func main() {
	u, _ := user.Current()

	runDir := "/home/" + u.Username + "/.sugoi"
	os.MkdirAll(runDir, 0755)
	os.Remove(path.Join(runDir, "audio.mp3"))
	os.Remove(path.Join(runDir, "audio.txt"))

	hasPipe := !term.IsTerminal(int(os.Stdin.Fd()))

	hasArg := len(os.Args) != 1

	if hasArg && !hasPipe {
		fmt.Println("No args passed!")
		return
	}

	arg := ""

	if hasArg {
		arg = os.Args[1]
	}

	var pipe []byte

	if hasPipe {
		pipe, _ = io.ReadAll(os.Stdin)
	}

	fmt.Println("Checking args....")

	result := ""

	if hasPipe {
		result = GenerateSummary(pipe)
	} else if IsURL(arg) {
		if IsSupported(arg) {
			fmt.Println("Transcribing " + color.YellowString(arg) + " with yt-dlp")
			// TODO: implement yt-dlp calling here

			DownloadAudio(arg, runDir)
			trans := Transcribe(arg, runDir)
			result = GenerateSummary(trans)
		} else {
			fmt.Println("[Error] URL was passed, but it's not supported by yt-dlp!")
			return
		}
	} else {
		info, stat := os.Stat(arg)

		// assume it's a string
		if os.IsNotExist(stat) {
			fmt.Println("Summarising string " + color.YellowString(`"`+arg+`"`))

			result = GenerateSummary([]byte(arg))
		} else if !info.IsDir() {
			// it's a file
			fmt.Println("Reading file " + color.YellowString(arg) + " and transcribing")

			data, err := os.ReadFile(arg)

			if err != nil {
				fmt.Println("Failed reading supplied file path contents")
				return
			}

			result = GenerateSummary(data)
		} else {
			fmt.Println(color.RedString("Directory paths are not supported!"))
			return
		}
	}

	fmt.Println("\n" + color.GreenString(result) + "\n")
}
