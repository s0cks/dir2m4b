package internal

import (
	"os"
	"path/filepath"
	"strings"
)

type AudioSourceVisitor func(file string) error

func isAudioSourceFile(file string) bool {
	return strings.HasSuffix(file, ".wav") || strings.HasSuffix(file, ".mp3")
}

func listAudioSourcesInDir(dir string) ([]string, error) {
	results := []string{}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		sourcePath := filepath.Join(dir, file.Name())
		if !isAudioSourceFile(sourcePath) {
			continue
		}

		results = append(results, sourcePath)
	}

	return results, nil
}

func VisitAudioSourcesInDir(dir string, vis AudioSourceVisitor) error {
	sources, err := listAudioSourcesInDir(dir)
	if err != nil {
		return nil
	}

	for _, sourcePath := range sources {
		if err := vis(sourcePath); err != nil {
			return err
		}
	}

	return nil
}

// package main
//
// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// )
//
// // Chapter represents an audiobook chapter with times in milliseconds
//
// func main() {
// 	// 1. Define your inputs and outputs
// 	inputFile := "input.mp3"
// 	outputFile := "output.m4b"
// 	coverFile := "cover.jpg"
// 	metaFile := "ffmpeg_meta.txt"
//
// 	// 3. Build the M4B using FFmpeg
// 	fmt.Println("Encoding M4B with chapters and cover art...")
//
// 	// Arguments:
// 	// -i input.mp3: Main audio input
// 	// -i ffmpeg_meta.txt: Metadata text file input
// 	// -i cover.jpg: Artwork input
// 	// -map_metadata 1: Use metadata from the text file (index 1)
// 	// -c:a aac -b:a 64k: Encode to standard audiobook AAC at 64kbps
// 	// -disposition:v attached_pic: Flags the cover image as artwork
// }
