package http

import (
	"context"
	"errors"
	json "github.com/bytedance/sonic"
	"log"
	"movie-app/metadata/pkg/model"
	"movie-app/movie/internal/v1/gateway"
	"movie-app/pkg/discovery"
	"net/http"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry *discovery.Registry) *Gateway {
	return &Gateway{registry: *registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addr, err := g.registry.GetRoundRobinAddress(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := addr + "/metadata/" + id
	log.Println("Calling metadata service. Request: GET" + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, errors.New("non-2xx response:" + resp.Status)
	}

	var md *model.Metadata
	if err := json.ConfigDefault.NewDecoder(resp.Body).Decode(md); err != nil {
		return nil, err
	}

	return md, nil
}
