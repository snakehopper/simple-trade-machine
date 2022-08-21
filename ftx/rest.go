package ftx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const URL = "https://ftx.com/api/"

func (a Api) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := a.sign(signaturePayload)
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", a.Api)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if a.Subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", a.Subaccount)
	}
	return req
}

func (a Api) _get(path string, body []byte) (*http.Response, error) {
	preparedRequest := a.signRequest("GET", path, body)
	resp, err := http.DefaultClient.Do(preparedRequest)
	return resp, err
}

func (a Api) _post(path string, body []byte) (*http.Response, error) {
	preparedRequest := a.signRequest("POST", path, body)
	resp, err := http.DefaultClient.Do(preparedRequest)
	return resp, err
}

func (a Api) _delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := a.signRequest("DELETE", path, body)
	resp, err := http.DefaultClient.Do(preparedRequest)
	return resp, err
}

func _processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error processing response: %v\n", err)
		return err
	}
	fmt.Println("<", strings.TrimSpace(string(body)))
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Printf("Error processing response: %v\n", err)
		return err
	}
	return nil
}
