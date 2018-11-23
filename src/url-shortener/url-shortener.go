package url_shortener

import (
	"encoding/json"
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
)

const (
	apiKey = "AIzaSyCrRYQY0HKi5B0ANQFcPmNJbYAY6Pb6aYs"
	apiUrl = "https://www.googleapis.com/urlshortener/v1/url"
)

type apiRequest struct {
	LongUrl string `json:"longUrl"`
}

type apiResponse struct {
	Kind    string `json:"kind"`
	Id      string `json:"id"`
	LongUrl string `json:"longUrl"`
}

func ShortenUrl(longUrl string) (string, error) {
	request := &apiRequest{
		LongUrl: longUrl,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	requestUrl := fmt.Sprintf("%s?key=%s", apiUrl, apiKey)
	bodyReader := bytes.NewReader(requestBody)
	rawResponse, err := http.Post(requestUrl, "application/json", bodyReader)
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
	fmt.Printf("JSON response: %s\n", responseBody)
	fmt.Printf("Parsed response: %+v\n", response)
	if err != nil {
		return "", err
	}

	fmt.Printf("Responding with: %+v\n", response)
	return response.Id, nil
}
