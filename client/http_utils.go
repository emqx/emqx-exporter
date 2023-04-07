package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func getHTTPClient(host string) *fasthttp.Client {
	return &fasthttp.Client{
		Name:                "EMQX-Exporter", //User-Agent
		MaxConnsPerHost:     5,
		MaxIdleConnDuration: 30 * time.Second,
		ReadTimeout:         3 * time.Second,
		WriteTimeout:        3 * time.Second,
		MaxConnWaitTimeout:  3 * time.Second,
		ConfigureClient: func(hc *fasthttp.HostClient) error {
			hc.Addr = host
			return nil
		},
	}
}

func callHTTPGet(client *fasthttp.Client, uri string) (data []byte, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(uri)
	req.SetHost("localhost")
	req.Header.SetMethod(http.MethodGet)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(*emqxUsername+":"+*emqxPassword)))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		err = fmt.Errorf("request %s failed. %w", uri, err)
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("%s: %s", uri, http.StatusText(resp.StatusCode()))
		return
	}

	//fmt.Fprintln(os.Stdout, "request:", uri, "code", resp.StatusCode(), "err", err, "body:", string(resp.Body()))
	return resp.Body(), nil
}

func callHTTPGetWithResp(client *fasthttp.Client, uri string, respData interface{}) (err error) {
	data, err := callHTTPGet(client, uri)
	if len(data) == 0 {
		err = fmt.Errorf("get nothing from api %s", uri)
		return
	}
	if !jsoniter.Valid(data) {
		err = errors.New("get response from api isn't valid json format")
		return
	}
	err = jsoniter.Unmarshal(data, respData)
	if err != nil {
		err = fmt.Errorf("unmarshal api resp failed: %s", uri)
		return
	}
	return
}
