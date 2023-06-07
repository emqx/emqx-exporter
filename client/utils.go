package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

var (
	readTimeout     = kingpin.Flag("emqx.readTimeout", "Maximum seconds for full response reading (including body)").Default("5").Int()
	connWaitTimeout = kingpin.Flag("emqx.connWaitTimeout", "Maximum seconds for waiting for a free connection").Default("5").Int()
)

func cutNodeName(nodeName string) string {
	slice := strings.Split(nodeName, "@")
	if len(slice) != 2 {
		return nodeName
	}

	if ip, err := netip.ParseAddr(slice[1]); err == nil {
		return ip.String()
	} else if pos := strings.IndexRune(slice[1], '.'); pos != 0 {
		return slice[1][:pos]
	}
	return slice[1]
}

func getHTTPClient(host string) *fasthttp.Client {
	return &fasthttp.Client{
		Name:                "EMQX-Exporter", //User-Agent
		MaxConnsPerHost:     5,
		MaxIdleConnDuration: 30 * time.Second,
		ReadTimeout:         time.Duration(*readTimeout) * time.Second,
		WriteTimeout:        5 * time.Second,
		MaxConnWaitTimeout:  time.Duration(*connWaitTimeout) * time.Second,
		ConfigureClient: func(hc *fasthttp.HostClient) error {
			hc.Addr = host
			return nil
		},
	}
}

func callHTTPGet(client *fasthttp.Client, uri string) (data []byte, statusCode int, err error) {
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

	statusCode = resp.StatusCode()

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("%s: %s", uri, http.StatusText(resp.StatusCode()))
		return
	}

	data = resp.Body()
	if len(data) == 0 {
		err = fmt.Errorf("get nothing from api %s", uri)
		return
	}
	if !jsoniter.Valid(data) {
		err = errors.New("get response from api isn't valid json format")
		return
	}

	errMsg := ""
	code := jsoniter.Get(data, "code")
	if code.ValueType() == jsoniter.NumberValue {
		// for emqx 4.4, it will return integer type code if occurred error
		if code.ToInt() != 0 {
			errMsg = fmt.Sprintf("%s: %d", uri, code.ToInt())
		}
	} else if code.ValueType() == jsoniter.StringValue {
		// for emqx 5, it will return string type code if occurred error
		if code.ToString() != "" {
			errMsg = fmt.Sprintf("%s: %s", uri, code.ToString())
		}
	}

	if errMsg != "" {
		msg := jsoniter.Get(data, "message")
		if msg.ValueType() == jsoniter.StringValue {
			errMsg = fmt.Sprintf("%s, msg=%s", errMsg, msg.ToString())
		}
		err = errors.New(errMsg)
	}
	return
}

func callHTTPGetWithResp(client *fasthttp.Client, uri string, respData interface{}) (err error) {
	data, _, err := callHTTPGet(client, uri)
	if err != nil {
		return
	}

	//fmt.Println("data:", string(data))
	err = jsoniter.Unmarshal(data, respData)
	if err != nil {
		err = fmt.Errorf("unmarshal api resp failed: %s, %s", uri, err.Error())
		return
	}
	return
}
