package data

import (
	"sync"
	"time"
)

type HostInfo struct {
	TelegramID    int64             `json:"telegramID"`
	HostID        string            `json:"hostID"`
	Nodes         []string          `json:"nodes"`
	Name          string            `json:"name"`
	Uptime        uint64            `json:"uptime"`
	Os            string            `json:"os"`
	OsVersion     string            `json:"osVersion"`
	KernelVersion string            `json:"kernelVersion"`
	Cpu           *CpuInfo          `json:"cpu"`
	Ram           uint64            `json:"ram"`
	Disks         map[string]uint64 `json:"disks"`
	IP            string            `json:"ip"`
}

type CpuInfo struct {
	Cores   int     `json:"cores"`
	Vendor  string  `json:"vendor"`
	Model   string  `json:"model"`
	Speed   float64 `json:"speed"`
	HasSSE4 bool    `json:"hasSSE4"`
}

type ResourcesUsage struct {
	TelegramID  int64  `json:"telegramID"`
	HostID      string `json:"hostID"`
	LastUpdated int64
	Cpu         float64            `json:"cpu"`
	Ram         float64            `json:"ram"`
	Disks       map[string]float64 `json:"disks"`
	DisksMut    sync.Mutex
}

type FullHostInfo struct {
	HostInfo       *HostInfo                 `json:"hostInfo"`
	ResourcesUsage *ResourcesUsage           `json:"resourcesUsage"`
	Nodes          map[string]*NodeLocalInfo `json:"nodes"`
	NodesMut       sync.Mutex
}

func (hi *HostInfo) GetName() string {
	if hi.Name != "" {
		return hi.Name
	}

	return hi.HostID
}

func (fhi *FullHostInfo) GetAllNodes() map[string]*NodeLocalInfo {
	res := make(map[string]*NodeLocalInfo)
	fhi.NodesMut.Lock()
	for k, v := range fhi.Nodes {
		res[k] = v
	}
	fhi.NodesMut.Unlock()

	return res
}

func (fhi *FullHostInfo) GetNodeByListenInterface(listenInterface string) *NodeLocalInfo {
	fhi.NodesMut.Lock()
	defer fhi.NodesMut.Unlock()

	return fhi.Nodes[listenInterface]
}

func (fhi *FullHostInfo) AddNode(node *NodeLocalInfo) {
	fhi.NodesMut.Lock()
	fhi.Nodes[node.ListenInterface] = node
	fhi.NodesMut.Unlock()
}

func (fhi *FullHostInfo) RemoveNode(node *NodeLocalInfo) {
	fhi.NodesMut.Lock()
	delete(fhi.Nodes, node.ListenInterface)
	fhi.NodesMut.Unlock()
}

func (fhi *FullHostInfo) GetNodeByKey(key string) *NodeLocalInfo {
	fhi.NodesMut.Lock()
	defer fhi.NodesMut.Unlock()

	for _, node := range fhi.Nodes {
		if node.Pubkey == key {
			return node
		}
	}

	return nil
}

func (fhi *FullHostInfo) UpdateNode(nodeInfo *NodeLocalInfo) {
	nodeInfo.LastUpdated = time.Now().Unix()
	fhi.NodesMut.Lock()
	fhi.Nodes[nodeInfo.ListenInterface] = nodeInfo
	fhi.NodesMut.Unlock()
}

func (fhi *FullHostInfo) UpdateResources(resInfo *ResourcesUsage) {
	resInfo.LastUpdated = time.Now().Unix()
	fhi.ResourcesUsage = resInfo
}
