package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"x-gwi/app/pki"
)

type HTTP struct {
	config     *ConfigClient
	httpClient *http.Client
}

func NewHTTP(config *ConfigClient, pki *pki.PKI) *HTTP {
	//nolint:exhaustruct,gomnd
	return &HTTP{
		config: config,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig:     pki.TLSConfigDial(),
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}
}

func (c *HTTP) Request(
	ctx context.Context, serviceMethod ServiceMethod, bodyRaw []byte) (
	[]byte, int, *http.Response, time.Duration, error) {
	t := time.Now()

	method, path, err := serviceMethod.HTTPMethodAndPath()
	if err != nil {
		return nil, 0, nil, time.Since(t), fmt.Errorf("error serviceMethod.HTTPMethodAndPath: %w", err)
	}

	u, err := url.Parse(fmt.Sprintf("https://%s", c.config.RESTGW.Address))
	if err != nil {
		return nil, 0, nil, time.Since(t), fmt.Errorf("error url.Parse: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.JoinPath(path).String(), bytes.NewReader(bodyRaw))
	if err != nil {
		return nil, 0, nil, time.Since(t), fmt.Errorf("error http.NewRequestWithContext: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, nil, time.Since(t), fmt.Errorf("error httpClient.Do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, res, time.Since(t), fmt.Errorf("error io.ReadAll: %w", err)
	}

	return body, len(body), res, time.Since(t), nil
}
