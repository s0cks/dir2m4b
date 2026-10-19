package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Chromaprint struct {
	Duration    float64 `json:"duration" yaml:"duration"`
	Fingerprint string  `json:"fingerprint" yaml:"fingerprint"`
}

type ManifestChapter struct {
	SourcePath string `json:"source" yaml:"source"`
	Hash       string `json:"hash" yaml:"hash"`
	Chromaprint
}

type Manifest struct {
	TargetPath string            `json:"target" yaml:"target"`
	Metadata   MetadataSpec      `json:"metadata" yaml:"metadata"`
	Hash       string            `json:"hash" yaml:"hash"`
	Timestamp  string            `json:"timestamp" yaml:"timestamp"`
	Root       string            `json:"root" yaml:"root"`
	Chapters   []ManifestChapter `json:"chapters" yaml:"chapters"`
	Chromaprint
}

func getChromaprint(source string) (*Chromaprint, error) {
	cmd := exec.Command("fpcalc", "-json", source)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("fpcalc failed: %v (stderr: %s)", err, stderr.String())
	}

	var result Chromaprint
	err = json.Unmarshal(stdout.Bytes(), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %v", err)
	}

	return &result, nil
}

func getHash(source string) (string, error) {
	file, err := os.Open(source)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to read file data: %w", err)
	}

	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

func CalculateMerkleRoot(hashes [][32]byte) ([32]byte, error) {
	if len(hashes) == 0 {
		return [32]byte{}, fmt.Errorf("cannot calculate Merkle root of an empty slice")
	}

	for len(hashes) > 1 {
		var nextLevel [][32]byte
		for i := 0; i < len(hashes); i += 2 {
			var left, right [32]byte
			left = hashes[i]

			if i+1 < len(hashes) {
				right = hashes[i+1]
			} else {
				right = left
			}

			concat := append(left[:], right[:]...)
			parentHash := sha256.Sum256(concat)
			nextLevel = append(nextLevel, parentHash)
		}

		hashes = nextLevel
	}

	return hashes[0], nil
}

func calcHashRoot(chapters []ManifestChapter) (string, error) {
	leaves := [][32]byte{}
	for _, ch := range chapters {
		leaves = append(leaves, sha256.Sum256([]byte(ch.Hash)))
	}

	root, err := CalculateMerkleRoot(leaves)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(root[:]), nil
}

func NewManifest(mdata *Metadata, outPath string) (*Manifest, error) {
	chapters := []ManifestChapter{}
	for _, ch := range mdata.Chapters {
		chromaprint, err := getChromaprint(ch.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to get chromaprint for %s: %v", ch.Source, err)
		}

		hash, err := getHash(ch.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to get hash for %s: %v", ch.Source, err)
		}

		chapters = append(chapters, ManifestChapter{
			SourcePath:  ch.Source[len(mdata.Root):],
			Chromaprint: *chromaprint,
			Hash:        hash,
		})
	}

	chromaprint, err := getChromaprint(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get chromaprint for %s: %v", outPath, err)
	}

	root, err := calcHashRoot(chapters)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hRoot: %v", err)
	}

	hash, err := getHash(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash for %s: %v", outPath, err)
	}

	manifest := Manifest{
		Chapters:    chapters,
		Root:        root,
		Metadata:    mdata.MetadataSpec,
		TargetPath:  outPath[strings.LastIndex(outPath, "/")+1:],
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Hash:        hash,
		Chromaprint: *chromaprint,
	}
	return &manifest, nil
}
