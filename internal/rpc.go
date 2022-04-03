package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURI = "https://api.codiy.net"
)

func Put(data map[string]string) error {
	httpClient := &http.Client{Timeout: time.Minute}
	httpContext := context.TODO()
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT", baseURI+"/api/tools/vocab", bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req = req.WithContext(httpContext)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	return nil
}
