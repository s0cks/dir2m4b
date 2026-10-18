package internal

// package main
//
// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// )
//
// // Chapter represents an audiobook chapter with times in milliseconds
// type Chapter struct {
// 	Title string
// 	Start int64 // Milliseconds
// 	End   int64 // Milliseconds
// }
//
// func main() {
// 	// 1. Define your inputs and outputs
// 	inputFile := "input.mp3"
// 	outputFile := "output.m4b"
// 	coverFile := "cover.jpg"
// 	metaFile := "ffmpeg_meta.txt"
//
// 	bookTitle := "My Great Audiobook"
// 	author := "Jane Doe"
//
// 	chapters := []Chapter{
// 		{Title: "Introduction", Start: 0, End: 60000},       // 0 to 1 min
// 		{Title: "Chapter 1", Start: 60000, End: 300000},     // 1 to 5 min
// 		{Title: "Chapter 2", Start: 300000, End: 600000},    // 5 to 10 min
// 	}
//
// 	// 2. Generate the FFmpeg metadata file
// 	err := writeMetadataFile(metaFile, bookTitle, author, chapters)
// 	if err != nil {
// 		fmt.Printf("Error writing metadata file: %v\n", err)
// 		return
// 	}
// 	defer os.Remove(metaFile) // Clean up the temp metadata text file when done
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
// 	cmd := exec.Command("ffmpeg",
// 		"-y",                 // Overwrite output file if it exists
// 		"-i", inputFile,
// 		"-i", metaFile,
// 		"-i", coverFile,
// 		"-map", "0:a",        // Take audio from the first input
// 		"-map", "2:v",        // Take video/image from the third input
// 		"-map_metadata", "1",
// 		"-c:a", "aac",
// 		"-b:a", "64k",
// 		"-c:v", "mjpeg",
// 		"-disposition:v", "attached_pic",
// 		outputFile,
// 	)
//
// 	// Capture output for debugging
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
//
// 	if err := cmd.Run(); err != nil {
// 		fmt.Printf("FFmpeg execution failed: %v\n", err)
// 		return
// 	}
//
// 	fmt.Printf("Successfully created %s!\n", outputFile)
// }
//
// // writeMetadataFile formats and saves the FFMETADATA1 configuration file
// func writeMetadataFile(filename string, title string, artist string, chapters []Chapter) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	// Write global headers
// 	_, _ = file.WriteString(";FFMETADATA1\n")
// 	_, _ = file.WriteString(fmt.Sprintf("title=%s\n", title))
// 	_, _ = file.WriteString(fmt.Sprintf("artist=%s\n", artist))
// 	_, _ = file.WriteString("genre=Audiobook\n\n")
//
// 	// Write individual chapters
// 	for _, ch := range chapters {
// 		_, _ = file.WriteString("[CHAPTER]\n")
// 		_, _ = file.WriteString("TIMEBASE=1/1000\n") // Millisecond timebase
// 		_, _ = file.WriteString(fmt.Sprintf("START=%d\n", ch.Start))
// 		_, _ = file.WriteString(fmt.Sprintf("END=%d\n", ch.End))
// 		_, _ = file.WriteString(fmt.Sprintf("title=%s\n\n", ch.Title))
// 	}
//
// 	return nil
// }
