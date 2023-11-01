package client

import (
	"emqx-exporter/config"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var ()

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

func getURI(metrics *config.Metrics) *fasthttp.URI {
	uri := &fasthttp.URI{}
	uri.SetUsername(metrics.APIKey)
	uri.SetPassword(metrics.APISecret)
	uri.SetScheme(metrics.Scheme)
	uri.SetHost(metrics.Target)
	return uri
}

func getHTTPClient(metrics *config.Metrics) *fasthttp.Client {
	return &fasthttp.Client{
		Name:                "EMQX-Exporter", //User-Agent
		MaxConnsPerHost:     5,
		MaxIdleConnDuration: 30 * time.Second,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        5 * time.Second,
		MaxConnWaitTimeout:  5 * time.Second,
		TLSConfig:           metrics.TLSClientConfig.ToTLSConfig(),
	}
}

func callHTTPGet(client *fasthttp.Client, uri *fasthttp.URI, requestURI string) (data []byte, statusCode int, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetURI(uri)
	req.URI().SetPath(requestURI)
	req.Header.SetMethod(http.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		err = fmt.Errorf("request %s failed. %w", req.URI().String(), err)
		return
	}

	statusCode = resp.StatusCode()

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("%s: %s", req.URI().String(), http.StatusText(resp.StatusCode()))
		return
	}

	data = resp.Body()
	if len(data) == 0 {
		err = fmt.Errorf("get nothing from api %s", req.URI().String())
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
			errMsg = fmt.Sprintf("%s: %d", req.URI().String(), code.ToInt())
		}
	} else if code.ValueType() == jsoniter.StringValue {
		// for emqx 5, it will return string type code if occurred error
		if code.ToString() != "" {
			errMsg = fmt.Sprintf("%s: %s", req.URI().String(), code.ToString())
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

func callHTTPGetWithResp(client *fasthttp.Client, uri *fasthttp.URI, requestURI string, respData interface{}) (err error) {
	data, _, err := callHTTPGet(client, uri, requestURI)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(data, respData)
	if err != nil {
		err = fmt.Errorf("unmarshal api resp failed: %s, %s", requestURI, err.Error())
		return
	}
	return
}
