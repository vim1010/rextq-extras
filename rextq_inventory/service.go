package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Service struct {
	Client  *http.Client
	BaseURL string
	User    string
	Pass    string
}

func NewService(url string, user string, pass string) *Service {
	return &Service{
		Client:  &http.Client{},
		BaseURL: url,
		User:    user,
		Pass:    pass,
	}
}

func (r *Service) Call(route string, opts map[string]any) ([]map[string]any, error) {
	var t string
	res := make([]map[string]any, 0)
	t = "{}"
	if opts != nil {
		tp, err := json.Marshal(opts)
		if err != nil {
			return res, errors.New("cannot convert opts to json")
		}
		t = string(tp)
	}
	url := fmt.Sprintf("%s/rpc/%s", r.BaseURL, route)
	d := bytes.NewBuffer([]byte(t))
	req, err := http.NewRequest("POST", url, d)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(r.User, r.Pass)
	resp, err := r.Client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	if resp.StatusCode != 200 {
		return res, errors.New(string(body))
	}
	err = json.Unmarshal(body, &res)
	return res, err
}

func (r *Service) First(route string, opts map[string]any) (row map[string]any, code int, err error) {
	rows := make([]map[string]any, 0)
	rows, err = r.Call(route, opts)
	if err != nil {
		return row, code, err
	}
	if len(rows) == 0 {
		return row, code, err
	}
	row = rows[0]
	return row, code, err
}
