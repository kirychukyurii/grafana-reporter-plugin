package grafana

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/apperrors"
	"net/http"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
)

// Client is a Grafana API client.
type Client struct {
	setting    *config.ReporterAppConfig
	connection *fasthttp.Client
}

func New(setting *config.ReporterAppConfig) (*Client, error) {
	cli := fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: setting.GrafanaConfig.InsecureSkipVerify,
		},
	}

	return &Client{
		setting:    setting,
		connection: &cli,
	}, nil
}

func (c *Client) Request(ctx context.Context, requestMethod, requestUrl string, requestPayload any, responseStruct any) error {
	var (
		err  error
		resp *fasthttp.Response
		body []byte
	)

	retryStatusCodes := c.setting.GrafanaConfig.RetryStatusCodesArr()
	if len(retryStatusCodes) == 0 {
		retryStatusCodes = []string{"429", "5xx"}
	}

	for n := 0; n <= c.setting.GrafanaConfig.RetryNum; n++ {

		// wait a bit if that's not the first request
		if n != 0 {
			if c.setting.GrafanaConfig.RetryTimeout == 0 {
				c.setting.GrafanaConfig.RetryTimeout = time.Second * 5
			}
			time.Sleep(c.setting.GrafanaConfig.RetryTimeout)
		}

		// If err is not nil, retry again
		// That's either caused by client policy, or failure to speak HTTP (such as network connectivity problem). A
		// non-2xx status code doesn't cause an error.
		resp, err = c.newRequest(requestMethod, requestUrl, requestPayload)
		if err != nil {
			continue
		}

		shouldRetry, err := matchRetryCode(resp.StatusCode(), retryStatusCodes)
		if err != nil {
			return err
		}

		if !shouldRetry {
			break
		}
	}

	if err != nil {
		return err
	}

	// do we need to decompress the response?
	contentEncoding := resp.Header.Peek("Content-Encoding")
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		body, err = resp.BodyGunzip()
		if err != nil {
			return fmt.Errorf("decompress the response: %v", err)
		}
	} else {
		body = resp.Body()
	}

	switch {
	case resp.StatusCode() == http.StatusNotFound:
		return fmt.Errorf("%v, body: %s", apperrors.ErrObjectNotFound, string(body))

	case resp.StatusCode() >= 400:
		return fmt.Errorf("expected status code %d but got %d, response: %s", fasthttp.StatusOK, resp.StatusCode(), string(body))
	}

	if responseStruct == nil {
		return nil
	}

	if err = json.Unmarshal(body, responseStruct); err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(requestMethod, requestUrl string, requestPayload any) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(requestUrl)
	req.Header.SetMethod(requestMethod)
	req.Header.Set("Accept", "application/json")

	if requestPayload != nil {
		payload, ok := requestPayload.([]byte)
		if ok {
			req.Header.SetContentType("application/json")
			req.SetBody(payload)
		}
	}

	if auth := c.setting.GrafanaConfig.BasicAuth(); auth != "" {
		req.Header.Add("Authorization", auth)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.connection.Do(req, resp); err != nil {
		return nil, fmt.Errorf("client get failed: %v", err)
	}

	return resp, nil
}

// matchRetryCode checks if the status code matches any of the configured retry status codes.
func matchRetryCode(gottenCode int, retryCodes []string) (bool, error) {
	gottenCodeStr := strconv.Itoa(gottenCode)

	for _, retryCode := range retryCodes {
		matched := true
		for i := range retryCode {
			c := retryCode[i]
			if c == 'x' {
				continue
			}

			if gottenCodeStr[i] != c {
				matched = false
				break
			}
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}
