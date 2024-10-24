package http

import (
	"context"
	"encoding/json"
	"fmt"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/rating/pkg/model"
	"net/http"
)

type Gatway struct {
	addr string
}

func New(addr string) *Gatway {
	return &Gatway{addr}
}

func (g *Gatway) GetAggregatedRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}
	var v float64
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return 0, err
	}
	return v, nil

}

func (g *Gatway) PutRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	req, err := http.NewRequest(http.MethodPut, g.addr+"/rating", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}
	return nil

}
