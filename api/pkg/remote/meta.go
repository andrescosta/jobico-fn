package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrescosta/goico/pkg/env"
)

type MetadataClient struct {
}

func NewMetadataClient() *MetadataClient {
	return &MetadataClient{}
}

func (c *MetadataClient) GetMetadata(ctx context.Context, name string) (map[string]string, error) {
	host := env.Env(name + ".host")
	url := fmt.Sprintf("http://%s/%s", host, "meta")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
