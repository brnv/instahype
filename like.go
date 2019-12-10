package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const setLikeURLFormat = "https://www.instagram.com/web/likes/%s/like/"

type LikeResponse struct {
	Status string `json:"status"`
}

func setLike(mediaID string, sessionID string) error {
	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf(setLikeURLFormat, mediaID),
		nil,
	)
	if err != nil {
		return err
	}

	request.AddCookie(&http.Cookie{
		Name:  "sessionid",
		Value: sessionID,
	})

	request.Header.Set("X-CSRFToken", "missed")

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("can't make request: %s", err.Error())
	}

	responseRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("can't read response: %s", err.Error())
	}

	responseData := LikeResponse{}

	err = json.Unmarshal(responseRaw, &responseData)
	if err != nil {
		return fmt.Errorf("can't unmarshal response body: %s", err.Error())
	}

	logger.Debugf("status: %s", responseData.Status)

	return nil
}
