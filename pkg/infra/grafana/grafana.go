package grafana

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/apperrors"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
)

// Client is a Grafana API client.
type Client struct {
	setting    *config.ReporterAppConfig
	baseURL    *url.URL
	connection *fasthttp.Client
}

func New(setting *config.ReporterAppConfig) (*Client, error) {
	u, err := url.Parse(setting.GrafanaConfig.URL)
	if err != nil {
		return nil, err
	}

	cli := fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: setting.GrafanaConfig.InsecureSkipVerify,
		},
	}

	return &Client{
		setting:    setting,
		baseURL:    u,
		connection: &cli,
	}, nil
}

func (c *Client) Request(ctx context.Context, requestMethod, requestPath string, requestPayload any, responseStruct any) error {
	var (
		err         error
		shouldRetry bool
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

			time.Sleep(time.Second * c.setting.GrafanaConfig.RetryTimeout)
		}

		// If err is not nil, retry again
		// That's either caused by client policy, or failure to speak HTTP (such as network connectivity problem). A
		// non-2xx status code doesn't cause an error.
		shouldRetry, err = c.newRequest(retryStatusCodes, requestMethod, requestPath, requestPayload, responseStruct)
		if !shouldRetry {
			break
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(retryStatusCodes []string, requestMethod, requestPath string, requestPayload any, responseStruct any) (bool, error) {
	var (
		err  error
		body []byte
	)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	u := c.baseURL
	u.Path = path.Join(u.Path, requestPath)
	req.SetRequestURI(u.String())
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

	if err = c.connection.Do(req, resp); err != nil {
		return false, fmt.Errorf("client %s get failed: %v", req.RequestURI(), err)
	}

	shouldRetry, err := matchRetryCode(resp.StatusCode(), retryStatusCodes)
	if err != nil {
		return false, err
	}

	if !shouldRetry {

		// do we need to decompress the response?
		contentEncoding := resp.Header.Peek("Content-Encoding")
		if bytes.EqualFold(contentEncoding, []byte("gzip")) {
			body, err = resp.BodyGunzip()
			if err != nil {
				return false, fmt.Errorf("decompress the response: %v", err)
			}
		} else {
			body = resp.Body()
		}

		switch {
		case resp.StatusCode() == http.StatusNotFound:
			return false, fmt.Errorf("%v, body: %s", apperrors.ErrObjectNotFound, string(body))

		case resp.StatusCode() >= 400:
			return false, fmt.Errorf("expected status code %d but got %d, response: %s", fasthttp.StatusOK, resp.StatusCode(), string(body))
		}

		if responseStruct == nil {
			return false, nil
		}

		if err = json.Unmarshal(body, &responseStruct); err != nil {
			return false, fmt.Errorf("unmarshal json (%s): %v", body, err)
		}
	}

	return shouldRetry, nil
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
