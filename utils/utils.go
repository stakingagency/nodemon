package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var (
	AppVersion string
	log        = logger.GetOrCreate("nodemon-utils")
)

const (
	NODEMON_GITHUB_REPO = "github.com/stakingagency/nodemon/cmd/nodemon"

	LISTEN_NODESMON_ROOT    = "/nodesmon"
	LISTEN_HOST_INFO        = LISTEN_NODESMON_ROOT + "/hostInfo"
	LISTEN_HOST_RESOURCES   = LISTEN_NODESMON_ROOT + "/hostResources"
	LISTEN_NODE_INFO        = LISTEN_NODESMON_ROOT + "/nodeInfo"
	LISTEN_HOST_TASKS       = LISTEN_NODESMON_ROOT + "/getTasks"
	LISTEN_HOST_TASK_RESULT = LISTEN_NODESMON_ROOT + "/sendTaskResult"

	HOST_CMD_REBOOT     = "reboot"
	HOST_CMD_UPDATE_APP = "updateApp"
	HOST_CMD_UPDATE_OS  = "updateOS"
	HOST_CMD_EXEC       = "exec"
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
