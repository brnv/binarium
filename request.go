package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RequestResponse struct {
	Data   Transaction `json:"data"`
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func request(jwt string, form url.Values) error {
	request, err := http.NewRequest(
		"POST",
		"https://api.binarium.com/api/v1/users/self/options",
		strings.NewReader(form.Encode()),
	)

	if err != nil {
		return fmt.Errorf("can't init request: %s", err.Error())
	}

	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "en-US")
	request.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Origin", "https://binarium.com")
	request.Header.Set("Referer", "https://binarium.com/terminal")
	request.Header.Set("Sec-Fetch-Dest", "empty")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Site", "same-site")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	request.Header.Set("X-JWT", jwt)

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("can't make request: %s", err.Error())
	}

	responseRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	responseBody := decompressIfGzipped(responseRaw)

	result := RequestResponse{}

	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		logger.Errorf("response body: %s", string(responseBody))
		return fmt.Errorf("can't unmarshal: %s", err.Error())
	}

	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			logger.Infof("response error code: '%s', message: '%s'", err.Code, err.Message)
		}
	}

	if result.Data.ID != 0 {
		logger.Infof("transaction id: %d", result.Data.ID)
	}

	return nil
}
