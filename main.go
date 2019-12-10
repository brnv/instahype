package main

import (
	"math/rand"
	"time"

	"github.com/kovetskiy/godocs"
	"github.com/kovetskiy/lorg"
)

var (
	logger  = lorg.NewLog()
	version = "[manual build]"
)

const usage = `name

Usage:
    name [options]
    name -h | --help

Options:
    --debug    Enable debug output.
    --trace    Enable trace output.
    -h --help  Show this help.
`

func main() {
	args := godocs.MustParse(usage, version, godocs.UsePager)

	logger.SetIndentLines(true)

	if args["--debug"].(bool) {
		logger.SetLevel(lorg.LevelDebug)
	}

	if args["--trace"].(bool) {
		logger.SetLevel(lorg.LevelTrace)
	}

	var (
		hashtag   = ""
		sessionID = ""
	)

	rand.Seed(time.Now().UnixNano())

	goalCount := rand.Intn(20) + 90

	count := 0

	for {
		if count == goalCount {
			break
		}

		logger.Infof("get videos by tag '%s'", hashtag)

		videos, err := getVideos(hashtag, 10, sessionID)
		if err != nil {
			logger.Fatalf("can't get videos: %s", err.Error())
		}

		logger.Debugf("videos count: %d", len(videos))

		for _, video := range videos {
			logger.Debugf(
				"set like: https://www.instagram.com/p/%s", video.Shortcode,
			)

			err = setLike(video.ID, sessionID)
			if err != nil {
				logger.Errorf(
					"media id: %s, error: %s",
					video.Shortcode, err.Error(),
				)
				continue
			}

			count++

			logger.Infof("progress %d/%d", count, goalCount)

			if count == goalCount {
				break
			}

			waitInterval := time.Second * (time.Duration)(rand.Intn(7)+14)

			logger.Debugf("sleep delay: %s", waitInterval.String())

			time.Sleep(waitInterval)
		}
	}
}
