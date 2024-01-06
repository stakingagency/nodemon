package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	AppVersion string
)

const (
	LISTEN_NODESMON_ROOT  = "/nodesmon"
	LISTEN_HOST_INFO      = LISTEN_NODESMON_ROOT + "/hostInfo"
	LISTEN_HOST_RESOURCES = LISTEN_NODESMON_ROOT + "/hostResources"
	LISTEN_NODE_INFO      = LISTEN_NODESMON_ROOT + "/nodeInfo"
)

func PostHTTP(address, body string, timeout ...time.Duration) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, address, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	if len(timeout) == 0 {
		client.Timeout = time.Minute
	} else {
		client.Timeout = timeout[0]
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	return io.ReadAll(resp.Body)
}

func PostJsonHTTP(address string, body interface{}, timeout ...time.Duration) ([]byte, error) {
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return PostHTTP(address, string(bytes), timeout...)
}

func GetHTTP(address string, body string, timeout ...time.Duration) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, address, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	if len(timeout) == 0 {
		client.Timeout = time.Minute
	} else {
		client.Timeout = timeout[0]
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return resBody, fmt.Errorf("http error %v %v, endpoint %s", resp.StatusCode, resp.Status, address)
	}

	return resBody, nil
}
