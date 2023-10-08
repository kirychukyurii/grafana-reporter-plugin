package grafana

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"

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

func (c *Client) Request(ctx context.Context, requestMethod, requestUrl string, requestPayload any) ([]byte, error) {
	var err error

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

	if err = c.connection.Do(req, resp); err != nil {
		return nil, fmt.Errorf("client get failed: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("expected status code %d but got %d, response: %s", fasthttp.StatusOK, resp.StatusCode(), resp.Body())
	}

	// Do we need to decompress the response?
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		body, err = resp.BodyGunzip()
		if err != nil {
			return nil, fmt.Errorf("decompress the response: %v", err)
		}
	} else {
		body = resp.Body()
	}

	return body, nil
}
