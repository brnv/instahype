package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Media struct {
	ID        string `json:"id"`
	Shortcode string `json:"shortcode"`
	IsVideo   bool   `json:"is_video"`
}

type MediaResponse struct {
	Data struct {
		Hashtag struct {
			EdgeHashtagToMedia struct {
				Edges []struct {
					Node Media `json:"node"`
				} `json:"edges"`
			} `json:"edge_hashtag_to_media"`
		} `json:"hashtag"`
	} `json:"data"`
}

const (
	queryHash = "90cba7a4c91000cf16207e4f3bee2fa2"

	getMediaURLFormat = "https://www.instagram.com/graphql/query/?query_hash=" +
		queryHash +
		"&variables={\"tag_name\":\"%s\",\"first\":50}"
)

func getVideos(hashtag string, limit int, sessionID string) ([]Media, error) {
	request, err := http.NewRequest(
		"GET",
		fmt.Sprintf(getMediaURLFormat, hashtag),
		nil,
	)
	if err != nil {
		return []Media{}, err
	}

	request.AddCookie(&http.Cookie{
		Name:  "sessionid",
		Value: sessionID,
	})

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return []Media{}, fmt.Errorf("can't make request: %s", err.Error())
	}

	responseRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []Media{}, fmt.Errorf("can't read response: %s", err.Error())
	}

	responseData := MediaResponse{}

	err = json.Unmarshal(responseRaw, &responseData)
	if err != nil {
		return []Media{}, fmt.Errorf(
			"can't unmarshal response body: %s", err.Error(),
		)
	}

	videos := []Media{}

	count := 0

	for _, media := range responseData.Data.Hashtag.EdgeHashtagToMedia.Edges {
		if !media.Node.IsVideo {
			continue
		}

		videos = append(videos, media.Node)

		count++
		if count == limit {
			break
		}
	}

	return videos, nil
}
