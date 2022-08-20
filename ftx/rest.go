package ftx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const URL = "https://ftx.com/api/"

func (c *Client) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := c.sign(signaturePayload)
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", c.Api)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if c.Subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", c.Subaccount)
	}
	return req
}

func (c *Client) _get(path string, body []byte) (*http.Response, error) {
	preparedRequest := c.signRequest("GET", path, body)
	resp, err := c.Client.Do(preparedRequest)
	return resp, err
}

func (c *Client) _post(path string, body []byte) (*http.Response, error) {
	preparedRequest := c.signRequest("POST", path, body)
	resp, err := c.Client.Do(preparedRequest)
	return resp, err
}

func (c *Client) _delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := c.signRequest("DELETE", path, body)
	resp, err := c.Client.Do(preparedRequest)
	return resp, err
}

func _processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error processing response: %v\n", err)
		return err
	}
	fmt.Println("<", string(body))
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Printf("Error processing response: %v\n", err)
		return err
	}
	return nil
}
