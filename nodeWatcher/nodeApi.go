package nodeWatcher

import (
	"encoding/json"
	"strings"

	"github.com/stakingagency/nodemon/data"
	"github.com/stakingagency/nodemon/utils"
)

func (nw *NodeWatcher) queryNode(endpoint string) ([]byte, error) {
	endpoint = "http://" + strings.ReplaceAll(nw.listenInterface+endpoint, "//", "/")

	return utils.GetHTTP(endpoint, "")
}

func (nw *NodeWatcher) getNodeStatus() (*data.NodeApiStatus, error) {
	bytes, err := nw.queryNode("/node/status")
	if err != nil {
		return nil, err
	}

	res := &data.NodeApiStatusResponse{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	return res.Data.Metrics, nil
}

func (nw *NodeWatcher) getNodeP2PStatus() (*data.NodeApiP2PStatus, error) {
	bytes, err := nw.queryNode("/node/p2pstatus")
	if err != nil {
		return nil, err
	}

	res := &data.NodeApiP2PStatusResponse{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	return res.Data.Metrics, nil
}

func (nw *NodeWatcher) getNodeBootstrapStatus() (*data.NodeApiBootstrapStatus, error) {
	bytes, err := nw.queryNode("/node/bootstrapstatus")
	if err != nil {
		return nil, err
	}

	res := &data.NodeApiBootstrapStatusResponse{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	return res.Data.Metrics, nil
}

func (nw *NodeWatcher) getNodePeerInfo() (*data.NodeApiPeerInfo, error) {
	bytes, err := nw.queryNode("/node/peerinfo")
	if err != nil {
		return nil, err
	}

	res := &data.NodeApiPeerInfoResponse{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	return res.Data.Info, nil
}

func (nw *NodeWatcher) getNodeManagedKeys() ([]string, error) {
	bytes, err := nw.queryNode("/node/managed-keys")
	if err != nil {
		return nil, err
	}

	res := &data.NodeApiManagedKeysResponse{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	return res.Data.ManagedKeys, nil
}
