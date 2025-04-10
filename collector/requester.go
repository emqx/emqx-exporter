package collector

import (
	"emqx-exporter/config"
	"errors"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type requester struct {
	client *fasthttp.Client
	uri    *fasthttp.URI
}

func newRequester(metrics *config.Metrics) *requester {
	uri := &fasthttp.URI{}
	uri.SetUsername(metrics.APIKey)
	uri.SetPassword(metrics.APISecret)
	uri.SetScheme(metrics.Scheme)
	uri.SetHost(metrics.Target)

	return &requester{
		uri: uri,
		client: &fasthttp.Client{
			Name:                "EMQX-Exporter", //User-Agent
			MaxConnsPerHost:     5,
			MaxIdleConnDuration: 30 * time.Second,
			ReadTimeout:         5 * time.Second,
			WriteTimeout:        5 * time.Second,
			MaxConnWaitTimeout:  5 * time.Second,
			TLSConfig:           metrics.TLSClientConfig.ToTLSConfig(),
			DialDualStack:       true,
		},
	}
}

func (r *requester) callHTTPGet(requestURI string) (data []byte, statusCode int, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetURI(r.uri)
	req.URI().SetPath(requestURI)
	req.Header.SetMethod(http.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = r.client.Do(req, resp)
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

func (r *requester) callHTTPGetWithResp(requestURI string, respData interface{}) (err error) {
	data, _, err := r.callHTTPGet(requestURI)
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
