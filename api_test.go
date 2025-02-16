package main

import (
	"encoding/json"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	//"strings"
	"testing"

)


func Test_handler_notImplemented(t *testing.T) {
	// instantiate mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	// create a request to our mock HTTP server
	//    in our case it means to create DELETE request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	// check if the response from the handler is what we expect
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusNotImplemented, resp.StatusCode)
	}
}


func Test_handler_malformedURL(t *testing.T) {
	//the mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/igcinfo/987/",
		ts.URL + "/igc/",
	}
	for _, tstring := range testCases {
		resp, err := http.Get(tstring)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d", tstring,
				http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}



func Test_handler_getAllIds_empty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()


	resp, err := http.Get(ts.URL+"/igcinfo/api/igc")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	var a []interface{}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if len(a) != 0 {
		t.Errorf("Excpected empty array, got %s", a)
	}
}






