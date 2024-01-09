package data

type NodeApiStatusResponse struct {
	Data struct {
		Metrics *NodeApiStatus `json:"metrics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type NodeApiStatus struct {
	AccountsSnapshotInProgress uint64 `json:"erd_accounts_snapshot_in_progress"`
	AppVersion                 string `json:"erd_app_version"`
	AreVmQueriesReady          string `json:"erd_are_vm_queries_ready"`
	ChainId                    string `json:"erd_chain_id"`
	ConnectedNodes             uint64 `json:"erd_connected_nodes"`
	IntraShardValidatorNodes   uint64 `json:"erd_intra_shard_validator_nodes"`
	IsSyncing                  uint64 `json:"erd_is_syncing"`
	NetworkRecvBps             uint64 `json:"erd_network_recv_bps"`
	NetworkSentBps             uint64 `json:"erd_network_sent_bps"`
	DisplayName                string `json:"erd_node_display_name"`
	NodeType                   string `json:"erd_node_type"`
	Nonce                      uint64 `json:"erd_nonce"`
	ConnectedPeers             uint64 `json:"erd_num_connected_peers"`
	PeerType                   string `json:"erd_peer_type"`
	PeerSubType                string `json:"erd_peer_subtype"`
	PeersSnapshotInProgress    uint64 `json:"erd_peers_snapshot_in_progress"`
	Pubkey                     string `json:"erd_public_key_block_sign"`
	RedundancyIsMainActive     string `json:"erd_redundancy_is_main_active"`
	RedundancyLevel            string `json:"erd_redundancy_level"`
	ShardID                    uint32 `json:"erd_shard_id"`
	SynchronizedRound          uint64 `json:"erd_synchronized_round"`
	TxPoolLoad                 uint64 `json:"erd_tx_pool_load"`
}

type NodeApiP2PStatusResponse struct {
	Data struct {
		Metrics *NodeApiP2PStatus `json:"metrics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type NodeApiP2PStatus struct {
	CrossShardObservers             string `json:"erd_p2p_cross_shard_observers"`
	CrossShardObserversFullArchive  string `json:"erd_p2p_cross_shard_observers_full_archive"`
	CrossShardValidators            string `json:"erd_p2p_cross_shard_validators"`
	CrossShardValidatorsFullArchive string `json:"erd_p2p_cross_shard_validators_full_archive"`
	IntraShardObservers             string `json:"erd_p2p_intra_shard_observers"`
	IntraShardObserversFullArchive  string `json:"erd_p2p_intra_shard_observers_full_archive"`
	IntraShardValidators            string `json:"erd_p2p_intra_shard_validators"`
	IntraShardValidatorsFullArchive string `json:"erd_p2p_intra_shard_validators_full_archive"`
	PeerInfo                        string `json:"erd_p2p_peer_info"`
	UnknownShardPeers               string `json:"erd_p2p_unknown_shard_peers"`
}

type NodeApiPeerInfoResponse struct {
	Data struct {
		Info *NodeApiPeerInfo `json:"info"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type NodeApiPeerInfo struct {
	Isblacklisted bool     `json:"isblacklisted"`
	Pid           string   `json:"pid"`
	Pubkey        string   `json:"pk"`
	PeerType      string   `json:"peertype"`
	PeerSubType   string   `json:"peersubtype"`
	Addresses     []string `json:"addresses"`
}

type NodeApiBootstrapStatusResponse struct {
	Data struct {
		Metrics *NodeApiBootstrapStatus `json:"metrics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type NodeApiBootstrapStatus struct {
	TrieSyncNumBytesReceived  uint64 `json:"erd_trie_sync_num_bytes_received"`
	TrieSyncNumNodesProcessed uint64 `json:"erd_trie_sync_num_nodes_processed"`
}

type NodeApiManagedKeysResponse struct {
	Data struct {
		ManagedKeys []string `json:"managedKeys"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type NodeLocalInfo struct {
	TelegramID      int64  `json:"telegramID"`
	HostID          string `json:"hostID"`
	ListenInterface string `json:"node"`
	LastUpdated     int64

	ChainId                    string `json:"erd_chain_id"`
	DisplayName                string `json:"erd_node_display_name"`
	Pubkey                     string `json:"erd_public_key_block_sign"`
	ShardID                    uint32 `json:"erd_shard_id"`
	NodeType                   string `json:"erd_node_type"`
	PeerType                   string `json:"erd_peer_type"`
	PeerSubType                string `json:"erd_peer_subtype"`
	Nonce                      uint64 `json:"erd_nonce"`
	AppVersion                 string `json:"erd_app_version"`
	ConnectedNodes             uint64 `json:"erd_connected_nodes"`
	ConnectedPeers             uint64 `json:"erd_num_connected_peers"`
	IsSyncing                  bool   `json:"erd_is_syncing"`
	SynchronizedRound          uint64 `json:"erd_synchronized_round"`
	RedundancyIsMainActive     bool   `json:"erd_redundancy_is_main_active"`
	RedundancyLevel            string `json:"erd_redundancy_level"`
	NetworkRecvBps             uint64 `json:"erd_network_recv_bps"`
	NetworkSentBps             uint64 `json:"erd_network_sent_bps"`
	AccountsSnapshotInProgress uint64 `json:"erd_accounts_snapshot_in_progress"`
	PeersSnapshotInProgress    uint64 `json:"erd_peers_snapshot_in_progress"`

	Prefs *Preferences `json:"preferences"`

	PeerInfo          string `json:"erd_p2p_peer_info"`
	UnknownShardPeers string `json:"erd_p2p_unknown_shard_peers"`

	TrieSyncNumBytesReceived  uint64 `json:"erd_trie_sync_num_bytes_received"`
	TrieSyncNumNodesProcessed uint64 `json:"erd_trie_sync_num_nodes_processed"`

	ManagedKeys []string `json:"managedKeys"`
}

type TaskResult struct {
	HostID string `json:"hostID"`
	Task   string `json:"command"`
	Output []byte `json:"output"`
	Error  string `json:"error"`
}
