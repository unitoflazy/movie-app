package http

import (
	"context"
	"errors"
	json "github.com/bytedance/sonic"
	"movie-app/metadata/pkg/model"
	"movie-app/movie/internal/v1/gateway"
	"net/http"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, errors.New("non-2xx response: " + resp.Status)
	}

	var md *model.Metadata
	if err := json.ConfigDefault.NewDecoder(resp.Body).Decode(md); err != nil {
		return nil, err
	}

	return md, nil
}
