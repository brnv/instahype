package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	loginDeviceID           = "android_AF59079D-5EB1-487A-84D0-196DDB107684"
	loginRequestURL         = "https://i.instagram.com/api/v1/accounts/login/"
	loginRequestContentType = "application/x-www-form-urlencoded"

	loginUserAgent = "Instagram 64.0.0.14.96 Android (23/6.0.1; 640dpi; 1440x2560; samsung; SM-G935F; hero2lte; samsungexynos8890; en_US; 125398467"
)

func login(username string, password string) (string, error) {
	formJson, err := json.Marshal(
		struct {
			Username string `json:"username"`
			Password string `json:"password"`
			DeviceID string `json:"device_id"`
		}{
			Username: username,
			Password: password,
			DeviceID: loginDeviceID,
		},
	)

	form := url.Values{
		"signed_body": []string{fmt.Sprintf("q.%s", string(formJson))},
	}

	request, err := http.NewRequest(
		"POST",
		loginRequestURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", loginRequestContentType)
	request.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))
	request.Header.Set("User-Agent", loginUserAgent)

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "sessionid" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("no session_id cookie found")
}
