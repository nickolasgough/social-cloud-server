package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Listener struct {
	Handler Handler
}

func (l *Listener) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	applyHeaders(&rw)
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Error(rw, err)
		return
	}
	fmt.Printf("Received request body: %s\n", body)

	request := l.Handler.Request()
	err = json.Unmarshal(body, request)
	if err != nil {
		l.Error(rw, err)
		fmt.Printf("Unmarshal errored with: %s\n", err.Error())
		return
	}
	fmt.Printf("Received request body: %+v\n", request)

	rawResponse, err := l.Handler.Process(context.Background(), request)
	if err != nil {
		l.Error(rw, err)
		fmt.Printf("Process errored with: %s\n", err.Error())
		return
	}
	fmt.Printf("Responding with: %+v\n", rawResponse)

	jsonResponse, err := json.Marshal(rawResponse)
	if err != nil {
		l.Error(rw, err)
		fmt.Printf("Marshal errored with: %s\n", err.Error())
		return
	}
	fmt.Printf("Responding with: %s\n", jsonResponse)

	fmt.Fprintf(rw, "%s", jsonResponse)
}

func (l *Listener) Error(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "Error: %s", err.Error())
}

func applyHeaders(rw *http.ResponseWriter) {
	(*rw).Header().Set("Content-Type", "application/json")
	(*rw).Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	(*rw).Header().Set("Access-Control-Allow-Methods", "POST, GET")
}
