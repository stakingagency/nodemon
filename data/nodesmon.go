package data

import (
	"sync"
)

type HostInfo struct {
	AppVersion    string            `json:"appVersion"`
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

func (hi *HostInfo) GetName() string {
	if hi.Name != "" {
		return hi.Name
	}

	return hi.HostID
}
