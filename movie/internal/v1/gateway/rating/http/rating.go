package http

import (
	"bytes"
	"context"
	"errors"
	json "github.com/bytedance/sonic"
	"movie-app/movie/internal/v1/gateway"
	"movie-app/rating/pkg/model"
	"net/http"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, ratingType model.RecordType, id model.RecordID) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(id))
	values.Add("type", string(ratingType))

	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, errors.New("non-2xx response: " + resp.Status)
	}

	var v float64
	if err := json.ConfigDefault.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

func (g *Gateway) PutRating(ctx context.Context, ratingType model.RecordType, id model.RecordID, rating *model.Rating) error {
	payload, err := json.Marshal(*rating)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, g.addr+"/rating", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	values := req.URL.Query()
	values.Add("id", string(id))
	values.Add("type", string(ratingType))

	req.URL.RawQuery = values.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return errors.New("non-2xx response: " + resp.Status)
	}

	return nil
}
