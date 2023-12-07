package remote

import (
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

func (c *MetadataClient) GetMetadata(name string) (*map[string]string, error) {
	host := env.GetAsString(name + ".host")
	url := fmt.Sprintf("http://%s/%s", host, "meta")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	r := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
