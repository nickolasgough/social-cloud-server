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
	rw.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Error(rw, err)
		return
	}

	request := l.Handler.Request()
	err = json.Unmarshal(b, request)
	if err != nil {
		l.Error(rw, err)
		return
	}

	rawResponse, err := l.Handler.Process(context.Background(), request)
	if err != nil {
		l.Error(rw, err)
		return
	}

	jsonResponse, err := json.Marshal(rawResponse)
	if err != nil {
		l.Error(rw, err)
		return
	}

	fmt.Fprintf(rw, "%s", jsonResponse)
}

func (l *Listener) Error(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "Error: %s", err.Error())
}
