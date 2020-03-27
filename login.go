package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type LoginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func login(email string, password string) (string, error) {
	form := url.Values{}

	form.Add("email", email)
	form.Add("password", password)
	form.Add("type", "jwt")

	request, err := http.NewRequest(
		"POST",
		"https://api.binarium.com/api/v1/login",
		strings.NewReader(form.Encode()),
	)

	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "en-US")
	request.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Origin", "https://binarium.com")
	request.Header.Set("Referer", "https://binarium.com/")
	request.Header.Set("Sec-Fetch-Dest", "empty")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Site", "same-site")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("can't make request: %s", err.Error())
	}

	responseRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	responseBody := decompressIfGzipped(responseRaw)

	result := LoginResponse{}

	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		logger.Errorf("response body: %s", string(responseBody))
		return "", fmt.Errorf("can't unmarshal: %s", err.Error())
	}

	if result.Data.Token != "" {
		err = saveToken(result.Data.Token)
		if err != nil {
			return result.Data.Token, err
		}
	}

	return result.Data.Token, nil
}

func saveToken(token string) error {
	path := "/home/operator/.cache/binarium"

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	filePath := path + "/token"

	ioutil.WriteFile(filePath, []byte(token), 0644)

	logger.Infof("token saved to file %s", filePath)

	return nil
}
