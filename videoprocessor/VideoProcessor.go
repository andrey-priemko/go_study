package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getVideoDuration(videoPath string) (float64, error) {
	result, err := exec.Command(`ffprobe`, `-v`, `error`, `-show_entries`, `format=duration`, `-of`, `default=noprint_wrappers=1:nokey=1`, videoPath).Output()
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(strings.Trim(string(result), "\n\r"), 64)
}

func ffmpegTimeFromSeconds(seconds int64) string {
	return time.Unix(seconds, 0).UTC().Format(`15:04:05.000000`)
}

func createVideoThumbnail(videoPath string, thumbnailPath string, thumbnailOffset int64) error {
	return exec.Command(`ffmpeg`, `-i`, videoPath, `-ss`, ffmpegTimeFromSeconds(thumbnailOffset), `-vframes`, `1`, thumbnailPath).Run()
}

func main() {
	duration, err := getVideoDuration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%f", duration)

	err = createVideoThumbnail(os.Args[1], os.Args[2], int64(duration) / 2)
	if err != nil {
		log.Fatal(err)
	}
}