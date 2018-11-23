package url_shortener

import (
	"encoding/json"
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
)

const (
	apiKey = "2e46a37bf73623ba49a12ecec034561470943189"
	apiUrl = "https://api-ssl.bitly.com/v4/bitlinks"
)

type apiRequest struct {
	LongUrl string `json:"long_url"`
}

type apiResponse struct {
	Link string `json:"link"`
}

func ShortenUrl(longUrl string) (string, error) {
	request := &apiRequest{
		LongUrl: longUrl,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	httpClient := &http.Client{}
	requestReader := bytes.NewReader(requestBody)
	postRequest, err := http.NewRequest("POST", apiUrl, requestReader)
	if err != nil {
		return "", err
	}
	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	rawResponse, err := httpClient.Do(postRequest)
	if err != nil {
		return "", err
	}

	defer rawResponse.Body.Close()
	responseBody, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return "", err
	}

	var response apiResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	return response.Link, nil
}
