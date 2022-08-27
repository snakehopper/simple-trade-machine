package com

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	log      *zap.SugaredLogger
	cli      *http.Client
	endpoint string
	apiKey   string
	secret   []byte
}

func (c Client) Get(path string, val url.Values, sign bool) (*http.Response, error) {
	ul, _ := url.Parse(c.endpoint)
	ul.Path = path
	ul.RawQuery = val.Encode()
	req, err := http.NewRequest("GET", ul.String(), nil)
	if err != nil {
		return nil, err
	}

	if sign {
		c.SignRequest(req)
	}

	return http.DefaultClient.Do(req)
}

func (c Client) Post(path string, val url.Values, sign bool) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprint(c.endpoint, path), strings.NewReader(val.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MBX-APIKEY", c.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if sign {
		c.SignRequest(req)
	}

	c.log.Info(">", SafeReadBody(req))
	return http.DefaultClient.Do(req)
}

func (c Client) Delete(path string, val url.Values, sign bool) (*http.Response, error) {
	ul, err := url.Parse(c.endpoint + path)
	if err != nil {
		return nil, err
	}
	ul.RawQuery = val.Encode()

	req, err := http.NewRequest("DELETE", ul.String(), strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MBX-APIKEY", c.apiKey)
	if sign {
		c.SignRequest(req)
	}

	c.log.Info(">", SafeReadBody(req))
	return http.DefaultClient.Do(req)
}

func (c Client) SignRequest(r *http.Request) {
	r.Header.Set("X-MBX-APIKEY", c.apiKey)

	h := c.HashRequest(r)
	v := r.URL.Query()
	v.Set("signature", h)
	r.URL.RawQuery = v.Encode()
}

func (c Client) HashRequest(r *http.Request) string {
	//inject timestamp epoch
	v := r.URL.Query()
	v.Set("timestamp", fmt.Sprint(time.Now().Unix()*1000))
	v.Set("recvWindow", "10000")
	r.URL.RawQuery = v.Encode()

	//start HASH jobs
	plain := r.URL.RawQuery + SafeReadBody(r)
	mac := hmac.New(sha256.New, c.secret)
	if _, err := mac.Write([]byte(plain)); err != nil {
		panic(err)
	}

	return hex.EncodeToString(mac.Sum(nil))
}

func SafeReadBody(r *http.Request) string {
	if r.Method != "PUT" && r.Method != "POST" {
		return ""
	}

	if r.GetBody != nil {
		body, err := r.GetBody()
		if err != nil {
			panic(err)
		}

		defer body.Close()

		bs, err := ioutil.ReadAll(body)
		if err != nil {
			panic(err)
		}
		return string(bs)
	}

	if r.Body == nil {
		return ""
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	//recover request content
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
	return string(bs)
}

func NewClient(logger *zap.SugaredLogger, endpoint, apiKey, secret string) *Client {
	return &Client{
		log:      logger,
		cli:      http.DefaultClient,
		endpoint: endpoint,
		apiKey:   apiKey,
		secret:   []byte(secret),
	}
}
