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

	responseJson, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	result := make(map[string]interface{})

	err = json.Unmarshal(responseJson, &result)
	if err != nil {
		logger.Errorf("responseJson: %s", string(responseJson))
		return fmt.Errorf("can't unmarshal: %s", err.Error())
	}

	logger.Info(result)

	return nil
}
