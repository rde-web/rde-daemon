package daemon

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"rde-daemon/internal/config"
)

type Client struct {
	service string
	client  http.Client
}

func NewClient(service string) *Client {
	return &Client{
		service: service,
		client: http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", path.Join(config.Instance.SocketsPath, service+".sock"))
				},
			},
		}}
}

func (c *Client) Post(route string, body io.Reader) ([]byte, error) {
	var uri url.URL = url.URL{
		Scheme: "http",
		Host:   "unix",
		Path:   route,
	}
	rsp, errDoPost := c.client.Post(uri.String(), "application/msgpack", body)
	if errDoPost != nil {
		return nil, errDoPost
	}
	defer rsp.Body.Close()
	if rsp.StatusCode < http.StatusBadRequest {
		return io.ReadAll(rsp.Body)
	}
	data, errReadBody := io.ReadAll(rsp.Body)
	if errReadBody != nil {
		return nil, fmt.Errorf(
			"bad status %d; cannot read resp body: %v",
			rsp.StatusCode,
			errReadBody,
		)
	}
	return nil, fmt.Errorf(
		"bad status %d; cannot read resp body: %v",
		rsp.StatusCode,
		string(data),
	)
}
