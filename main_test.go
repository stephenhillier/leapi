package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	api := &Server{nil, nil}
	handler := http.HandlerFunc(api.health)
	server := httptest.NewServer(handler)
	defer server.Close()

	response, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	received, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}

	expected := "OK"

	if expected != string(received) {
		t.Errorf("Expected '%s', received '%s'", expected, received)
	}

}
