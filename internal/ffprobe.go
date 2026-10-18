package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"charm.land/log/v2"
)

func GetDurationInMs(source string) (int64, error) {
	args := []string{
		"-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", source,
	}

	log.Debug(fmt.Sprintf("running %s %s", "ffprobe", strings.Join(args, " ")))

	cmd := exec.Command("ffprobe", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to run ffprobe: %v", err)
	}

	outStr := strings.TrimSpace(out.String())
	seconds, err := strconv.ParseFloat(outStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration from `%s`: %v", source, err)
	}

	return int64(seconds * 1000), nil
}
