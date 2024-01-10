package nodeWatcher

import (
	"time"

	"github.com/stakingagency/nodemon/data"
	"github.com/stakingagency/nodemon/utils"
)

func (nw *NodeWatcher) watchNode() {
	for !nw.terminate {
		time.Sleep(time.Second)
		status, err := nw.getNodeStatus()
		if err != nil {
			log.Error("get node status", "error", err)
			continue
		}

		p2pStatus, err := nw.getNodeP2PStatus()
		if err != nil {
			log.Error("get node p2p status", "error", err)
			continue
		}

		bootstrapStatus, err := nw.getNodeBootstrapStatus()
		if err != nil {
			log.Error("get node bootstrap status", "error", err)
			continue
		}

		managedKeys, err := nw.getNodeManagedKeys()
		if err != nil {
			log.Warn("get node managed keys", "error", err)
			continue
		}

		// peerInfo, err := nw.getNodePeerInfo()
		// if err != nil {
		// 	log.Error("get node peer info", "error", err)
		// 	continue
		// }

		localInfo := &data.NodeLocalInfo{
			TelegramID:      nw.appCfg.TelegramID,
			HostID:          nw.sysWatcher.GetHostID(),
			ListenInterface: nw.listenInterface,

			ChainId:                    status.ChainId,
			DisplayName:                status.DisplayName,
			Pubkey:                     status.Pubkey,
			ShardID:                    status.ShardID,
			NodeType:                   status.NodeType,
			PeerType:                   status.PeerType,
			PeerSubType:                status.PeerSubType,
			Nonce:                      status.Nonce,
			AppVersion:                 status.AppVersion,
			ConnectedNodes:             status.ConnectedNodes,
			ConnectedPeers:             status.ConnectedPeers,
			IsSyncing:                  status.IsSyncing != 0,
			SynchronizedRound:          status.SynchronizedRound,
			RedundancyIsMainActive:     status.RedundancyIsMainActive == "true",
			RedundancyLevel:            status.RedundancyLevel,
			NetworkRecvBps:             status.NetworkRecvBps,
			NetworkSentBps:             status.NetworkSentBps,
			AccountsSnapshotInProgress: status.AccountsSnapshotInProgress,
			PeersSnapshotInProgress:    status.PeersSnapshotInProgress,

			Prefs:             nw.prefs,
			ValidatorKey:      nw.validatorKey,
			AllValidatorsKeys: nw.allValidatorKeys,
			IsMultiKey:        nw.isMultiKey,

			PeerInfo:          p2pStatus.PeerInfo,
			UnknownShardPeers: p2pStatus.UnknownShardPeers,

			TrieSyncNumBytesReceived:  bootstrapStatus.TrieSyncNumBytesReceived,
			TrieSyncNumNodesProcessed: bootstrapStatus.TrieSyncNumNodesProcessed,

			ManagedKeys: managedKeys,
		}

		_, err = utils.PostJsonHTTP(nw.appCfg.Server+utils.LISTEN_NODE_INFO, localInfo)
		if err == nil {
			log.Info("sent node info", "interface", nw.listenInterface, "info", localInfo)
			time.Sleep(time.Minute)
		} else {
			log.Warn("send node info", "error", err)
		}
	}
}
