package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Listener struct {
	Handler Handler
}

func (l *Listener) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	go l.handle(rw, r)
}

func (l *Listener) handle(rw http.ResponseWriter, r *http.Request) {
	applyHeaders(&rw)
	defer r.Body.Close()

	fmt.Printf("Recieved request from %s of type %s\n", r.URL, r.Method)

	var body []byte
	var err error
	switch r.Method {
	case "GET":
		body, err = parseUrlParameters(r.URL.RawQuery)
		break
	case "POST":
		body, err = ioutil.ReadAll(r.Body)
		break
	}
	if err != nil {
		l.error(rw, err)
		return
	}
	fmt.Printf("Received request body: %s\n", body)

	request := l.Handler.Request()
	err = json.Unmarshal(body, request)
	if err != nil {
		l.error(rw, err)
		fmt.Printf("Unmarshal errored with: %s\n", err.Error())
		return
	}
	fmt.Printf("Received request body: %+v\n", request)

	rawResponse, err := l.Handler.Process(context.Background(), request)
	if err != nil {
		l.error(rw, err)
		fmt.Printf("Process errored with: %s\n", err.Error())
		return
	}
	fmt.Printf("Responding with: %+v\n", rawResponse)

	jsonResponse, err := json.Marshal(rawResponse)
	if err != nil {
		l.error(rw, err)
		fmt.Printf("Marshal errored with: %s\n", err.Error())
		return
	}
	fmt.Printf("Responding with: %s\n", jsonResponse)

	fmt.Fprintf(rw, "%s", jsonResponse)
}

func (l *Listener) error(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "Error: %s", err.Error())
}

func applyHeaders(rw *http.ResponseWriter) {
	(*rw).Header().Set("Content-Type", "application/json")
	(*rw).Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	(*rw).Header().Set("Access-Control-Allow-Methods", "POST, GET")
}

func parseUrlParameters(query string) ([]byte, error) {
	parameters := strings.Split(query, "&")
	var jsonParams []string
	for _, parameter := range parameters {
		keyValues := strings.Split(parameter, "=")
		jsonParam := fmt.Sprintf("\"%s\":\"%s\"", keyValues[0], keyValues[1])
		jsonParams = append(jsonParams, jsonParam)
	}
	return []byte(fmt.Sprintf("{%s}", strings.Join(jsonParams, ","))), nil
}
