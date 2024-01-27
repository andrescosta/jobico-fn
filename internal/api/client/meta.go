package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/rs/zerolog"
)

type Metadata struct{}

func NewMetadata() *Metadata {
	return &Metadata{}
}

func (c *Metadata) Metadata(ctx context.Context, name string) (map[string]string, error) {
	logger := zerolog.Ctx(ctx)
	host := env.String(name + ".host")
	url := fmt.Sprintf("http://%s/%s", host, "meta")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if req.Body != nil {
			if err := req.Body.Close(); err != nil {
				logger.Warn().Err(err)
			}
		}
	}()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				logger.Warn().Err(err)
			}
		}
	}()
	r := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
