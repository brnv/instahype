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

const usage = `instahype

Usage:
    instahype [options]
    instahype -h | --help

Options:
    --session-id <string>   Use specific session.
    --username <string>     Username to login to Instagram (if no session id specified).
    --password <string>     Password to login to Instagram (if no session id specified).
    --tag <string>          Hashtag to search for [default: guitar].
    --start-delay           Wait short delay before start work.
    --debug                 Enable debug output.
    --trace                 Enable trace output.
    --cookies-db <path>     Get session id from sqlite cookies database.
    -h --help               Show this help.
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
		err error

		sessionID string
		hashtag   = args["--tag"].(string)
	)

	rand.Seed(time.Now().UnixNano())

	if args["--start-delay"].(bool) {
		startDelay := time.Second * time.Duration(rand.Intn(600)+180)
		logger.Infof("wait for %s before work", startDelay)
		time.Sleep(startDelay)
	}

	if args["--cookies-db"] != nil {
		logger.Trace("use chromium session")
		sessionID, err = getSessionFromCookies(args["--cookies-db"].(string))
		if err != nil {
			logger.Fatal(err)
		}
	} else if args["--session-id"] != nil {
		logger.Trace("use provided session id")
		sessionID = args["--session-id"].(string)
	}

	if sessionID == "" {
		logger.Trace("login to Instagram")

		var (
			password string
			username string
		)

		if args["--password"] != nil {
			password = args["--password"].(string)
		}

		if args["--username"] != nil {
			username = args["--username"].(string)
		}

		sessionID, err = login(username, password)
		if err != nil {
			logger.Fatalf("can't login: %s", err.Error())
		}
	}

	logger.Debugf("hashtag: %s", hashtag)
	logger.Debugf("session id: %s", sessionID)

	rand.Seed(time.Now().UnixNano())

	goalCount := rand.Intn(10) + 40

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

			waitInterval := time.Second * time.Duration(rand.Intn(7)+14)

			logger.Debugf("sleep delay: %s", waitInterval.String())

			time.Sleep(waitInterval)
		}
	}
}
