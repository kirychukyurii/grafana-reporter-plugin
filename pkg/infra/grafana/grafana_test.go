package grafana

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
)

type mockServerCall struct {
	code int
	body string
}

type mockServer struct {
	upcomingCalls []mockServerCall
	executedCalls []mockServerCall
	server        *fasthttp.Server
}

func gapiMockServer(t *testing.T, code int, body string) (context.Context, *Client) {
	t.Helper()
	ctx := context.Background()

	return ctx, gapiMockServerFromCalls(t, []mockServerCall{{code, body}})
}

func gapiMockServerFromCalls(t *testing.T, calls []mockServerCall) *Client {
	t.Helper()

	mock := &mockServer{
		upcomingCalls: calls,
	}

	ln := fasthttputil.NewInmemoryListener()
	t.Cleanup(func() {
		ln.Close()
	})

	mock.server = &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			if len(mock.upcomingCalls) == 0 {
				t.Fatalf("unexpected call to %s %s", ctx.Request.Header.Method(), ctx.Request.RequestURI())
			}

			call := mock.upcomingCalls[0]
			if len(calls) > 1 {
				mock.upcomingCalls = mock.upcomingCalls[1:]
			} else {
				mock.upcomingCalls = nil
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(call.code)
			ctx.SetBodyString(call.body)

			mock.executedCalls = append(mock.executedCalls, call)
		},
	}

	go mock.server.Serve(ln)

	httpClient := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	client, err := New(&config.ReporterAppConfig{GrafanaConfig: config.GrafanaConfig{
		URL:              "http://grafana.foo.bar/test",
		RetryNum:         3,
		RetryTimeout:     5,
		RetryStatusCodes: "429",
		Username:         "admin",
		Password:         "admin",
	}})
	if err != nil {
		t.Fatal(err)
	}

	client.connection = httpClient

	return client
}

func TestGrafanaClient_InvalidURL(t *testing.T) {
	_, err := New(&config.ReporterAppConfig{GrafanaConfig: config.GrafanaConfig{
		URL: "://my-grafana.com", APIToken: "123",
	}})

	expected := "parse \"://my-grafana.com\": missing protocol scheme"
	if err.Error() != expected {
		t.Errorf("expected error: %v; got: %s", expected, err)
	}
}

func TestGrafanaClient_Response200(t *testing.T) {
	ctx, client := gapiMockServer(t, 200, `{"foo":"bar"}`)

	err := client.Request(ctx, "GET", "/foo", nil, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGrafanaClient_Response201(t *testing.T) {
	ctx, client := gapiMockServer(t, 201, `{"foo":"bar"}`)

	err := client.Request(ctx, "GET", "/foo", nil, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGrafanaClient_Response400(t *testing.T) {
	ctx, client := gapiMockServer(t, 400, `{"foo":"bar"}`)

	expected := `expected status code 200 but got 400, response: {"foo":"bar"}`
	err := client.Request(ctx, "GET", "/foo", nil, nil)
	if err != nil && err.Error() != expected {
		t.Errorf("expected error: %v; got: %s", expected, err)
	}
}

func TestGrafanaClient_Response500(t *testing.T) {
	ctx, client := gapiMockServer(t, 500, `{"foo":"bar"}`)

	expected := `expected status code 200 but got 500, response: {"foo":"bar"}`
	err := client.Request(ctx, "GET", "/foo", nil, nil)
	if err != nil && err.Error() != expected {
		t.Errorf("expected error: %v; got: %s", expected, err)
	}
}

func TestGrafanaClient_Response200Unmarshal(t *testing.T) {
	ctx, client := gapiMockServer(t, 200, `{"foo":"bar"}`)
	result := struct {
		Foo string `json:"foo"`
	}{}

	if err := client.Request(ctx, "GET", "/foo", nil, &result); err != nil {
		t.Fatal(err)
	}

	if result.Foo != "bar" {
		t.Errorf("expected: bar; got: %s", result.Foo)
	}
}

func TestGrafanaClient_RequestWithRetries(t *testing.T) {
	var try int

	ctx := context.Background()
	body := []byte(`lorem ipsum dolor sit amet`)

	ln := fasthttputil.NewInmemoryListener()
	t.Cleanup(func() {
		ln.Close()
	})

	// This is our actual test, checking that we do in fact receive a body.
	ts := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			try++

			got := ctx.Request.Body()
			if !bytes.Equal(body, got) {
				t.Errorf("retry %d: request body doesn't match body sent by client. exp: %v got: %v", try, body, got)
			}

			ctx.SetContentType("application/json")
			switch try {
			case 1:
				ctx.SetBodyString(`{"error":"waiting for the right time"}`)
				ctx.SetStatusCode(http.StatusInternalServerError)

			case 2:
				ctx.SetBodyString(`{"error":"calm down"}`)
				ctx.SetStatusCode(http.StatusTooManyRequests)

			case 3:
				ctx.SetBodyString(`{"foo":"bar"}`)

			default:
				t.Errorf("unexpected retry %d", try)
			}
		},
	}

	go ts.Serve(ln)

	httpClient := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	c, err := New(&config.ReporterAppConfig{GrafanaConfig: config.GrafanaConfig{
		URL:              "http://grafana.foo.bar/test",
		RetryNum:         3,
		RetryTimeout:     3,
		RetryStatusCodes: "429 500",
		Username:         "admin",
		Password:         "admin",
	}})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	c.connection = httpClient
	type res struct {
		Foo string `json:"foo"`
	}

	var got res
	if err = c.Request(ctx, http.MethodPost, "/", body, &got); err != nil {
		t.Fatalf("unexpected error sending request: %v", err)
	}

	exp := res{Foo: "bar"}
	if exp != got {
		t.Fatalf("response doesn't match. exp: %#v got: %#v", exp, got)
	}

	t.Logf("request successful after %d retries", try)
}
