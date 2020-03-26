package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type FilterResponse struct {
	Data []Transaction `json:"data"`
}

func findTransaction(jwt string, id int) (Transaction, error) {
	form := url.Values{}

	timeNow := time.Now().UTC()

	filterCreatedAtFrom := fmt.Sprintf(
		"%d-%02d-%02d",
		timeNow.Year(),
		timeNow.Month(),
		timeNow.Day()-1,
	)

	filterCreatedAtTo := fmt.Sprintf(
		"%d-%02d-%02dT23:59:59Z",
		timeNow.Year(),
		timeNow.Month(),
		timeNow.Day()+1,
	)

	form.Add("filter[type]", "2") // binary, not turbo
	form.Add("filter[asset]", "")
	form.Add("filter[createdAt][from]", filterCreatedAtFrom)
	form.Add("filter[createdAt][to]", filterCreatedAtTo)
	form.Add("filter[currency]", "1")

	printTransactionForm(form)

	request, err := http.NewRequest(
		"GET",
		"https://api.binarium.com/api/v1/users/self/options",
		nil,
	)

	if err != nil {
		return Transaction{}, fmt.Errorf("can't init request: %s", err.Error())
	}

	request.URL.RawQuery = form.Encode()

	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "en-US")
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Host", "api.binarium.com")
	request.Header.Set("Origin", "https://binarium.com")
	request.Header.Set("Referer", "https://binarium.com/main/deal-history")
	request.Header.Set("TE", "Trailers")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	request.Header.Set("X-JWT", jwt)

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return Transaction{}, fmt.Errorf("can't make request: %s", err.Error())
	}

	responseRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Transaction{}, err
	}

	responseBody := decompressIfGzipped(responseRaw)

	result := FilterResponse{}

	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		logger.Errorf("response body: %s", string(responseBody))
		return Transaction{}, fmt.Errorf("can't unmarshal: %s", err.Error())
	}

	resultTransaction := Transaction{}
	found := false
	for _, transation := range result.Data {
		if transation.ID == id {
			resultTransaction = transation
			found = true
			break
		}
	}

	if found {
		return resultTransaction, nil
	}

	return Transaction{}, fmt.Errorf("transaction %d not found", id)
}

func printTransactionForm(form url.Values) {
	info := fmt.Sprintf(
		"filter[type]: %s\n"+
			"filter[asset]: %s\n"+
			"filter[createdAt][from]: %s\n"+
			"filter[createdAt][to]: %s\n"+
			"filter[currency]: %s",
		form.Get("filter[type]"),
		form.Get("filter[asset]"),
		form.Get("filter[createdAt][from]"),
		form.Get("filter[createdAt][to]"),
		form.Get("filter[currency]"),
	)

	logger.Info(info)
}
