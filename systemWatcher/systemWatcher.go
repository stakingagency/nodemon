package systemWatcher

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/stakingagency/nodemon/data"
	"github.com/stakingagency/nodemon/utils"
)

type SystemWatcher struct {
	appCfg *data.NodeMonAppConfig

	hostInfo *data.HostInfo
	usage    *data.ResourcesUsage
}

var log = logger.GetOrCreate("system-watcher")

func NewSystemWatcher(appCfg *data.NodeMonAppConfig) (*SystemWatcher, error) {
	pc, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	flags := 0
	for _, flag := range cpuInfo[0].Flags {
		if flag == "sse4_1" || flag == "sse4_2" {
			flags++
		}
	}
	hasSSE4 := flags == 2

	cpuLoad, err := load.Avg()
	if err != nil {
		return nil, err
	}

	load := cpuLoad.Load15 * 100 / float64(len(cpuInfo))

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	disks := make(map[string]uint64)
	disksUsage := make(map[string]float64)
	for _, partition := range partitions {
		if !strings.Contains(partition.Opts, "rw") ||
			strings.HasPrefix(partition.Mountpoint, "/boot") ||
			strings.HasPrefix(partition.Mountpoint, "/run") ||
			strings.HasPrefix(partition.Mountpoint, "/sys") ||
			strings.HasPrefix(partition.Mountpoint, "/dev") {
			continue
		}

		total, used, err := getDiskInfo(partition.Mountpoint)
		if err != nil || total == 0 {
			continue
		}

		disks[partition.Mountpoint] = total
		disksUsage[partition.Mountpoint] = used
		log.Debug("get disks", "device", partition.Device, "fsType", partition.Fstype, "mountPoint", partition.Mountpoint, "opts", partition.Opts, "size", total, "used", used)
	}

	s := &SystemWatcher{
		appCfg: appCfg,
		hostInfo: &data.HostInfo{
			AppVersion:    utils.AppVersion,
			TelegramID:    appCfg.TelegramID,
			HostID:        pc.HostID,
			Nodes:         make([]string, 0),
			Name:          pc.Hostname,
			Uptime:        pc.Uptime,
			Os:            pc.Platform,
			OsVersion:     pc.PlatformVersion,
			KernelVersion: pc.KernelVersion,
			Cpu: &data.CpuInfo{
				Cores:   len(cpuInfo),
				Vendor:  cpuInfo[0].VendorID,
				Model:   cpuInfo[0].ModelName,
				Speed:   cpuInfo[0].Mhz,
				HasSSE4: hasSSE4,
			},
			Ram:   vm.Total,
			Disks: disks,
		},
		usage: &data.ResourcesUsage{
			TelegramID: appCfg.TelegramID,
			HostID:     pc.HostID,
			Cpu:        load,
			Ram:        vm.UsedPercent,
			Disks:      disksUsage,
		},
	}

	return s, nil
}

func (sw *SystemWatcher) GetAppConfig() *data.NodeMonAppConfig {
	return sw.appCfg
}

func (sw *SystemWatcher) GetHostID() string {
	return sw.hostInfo.HostID
}

func getDiskInfo(path string) (uint64, float64, error) {
	stat, err := disk.Usage(path)
	if err != nil {
		return 0, 0, nil
	}

	return stat.Total, stat.UsedPercent, nil
}

func (sw *SystemWatcher) StartTasks() {
	go sw.watchResourcesAndGetTasks()
}

func (sw *SystemWatcher) watchResourcesAndGetTasks() {
	_, err := utils.PostJsonHTTP(sw.appCfg.Server+utils.LISTEN_HOST_INFO, sw.hostInfo)
	if err != nil {
		log.Error("send host info", "error", err)
		os.Exit(1)
	}

	log.Info("sent host info", "info", sw.hostInfo)
	shouldResendHost := false
	lastUpdateHost := time.Now().Unix()

	for {
		time.Sleep(time.Second)

		pc, err := host.Info()
		if err != nil {
			log.Error("get host info", "error", err)
			continue
		}

		sw.hostInfo.Uptime = pc.Uptime

		cpuLoad, err := load.Avg()
		if err != nil {
			log.Error("get cpu load", "error", err)
			continue
		}

		newLoad := cpuLoad.Load15 * 100 / float64(sw.hostInfo.Cpu.Cores)
		sw.usage.Cpu = newLoad

		vm, err := mem.VirtualMemory()
		if err != nil {
			log.Error("get memory info", "error", err)
			continue
		}

		sw.usage.Ram = vm.UsedPercent

		for path := range sw.hostInfo.Disks {
			_, used, err := getDiskInfo(path)
			if err != nil {
				log.Error("get disk info", "error", err, "path", path)
				continue
			}

			sw.usage.DisksMut.Lock()
			sw.usage.Disks[path] = used
			sw.usage.DisksMut.Unlock()
		}

		now := time.Now().Unix()
		if now-lastUpdateHost >= 300 {
			shouldResendHost = true
		}

		bytes, err := utils.PostJsonHTTP(sw.appCfg.Server+utils.LISTEN_HOST_RESOURCES, sw.usage)
		if err == nil {
			log.Info("sent resources usage", "usage", sw.usage)
			if shouldResendHost {
				_, err := utils.PostJsonHTTP(sw.appCfg.Server+utils.LISTEN_HOST_INFO, sw.hostInfo)
				if err == nil {
					shouldResendHost = false
					lastUpdateHost = now
				}
			}
			tasks := make([]string, 0)
			err = json.Unmarshal(bytes, &tasks)
			if err == nil && len(tasks) > 0 {
				sw.executeTasks(tasks)
			}
			time.Sleep(time.Minute)
		} else {
			shouldResendHost = true
			log.Warn("send resources usage", "error", err)
		}
	}
}

func (sw *SystemWatcher) AddNode(node string) {
	found := false
	for _, oldNode := range sw.hostInfo.Nodes {
		if node == oldNode {
			found = true
			break
		}
	}

	if !found {
		sw.hostInfo.Nodes = append(sw.hostInfo.Nodes, node)
	}
}

func (sw *SystemWatcher) RemoveNode(node string) {
	for i, oldNode := range sw.hostInfo.Nodes {
		if node == oldNode {
			sw.hostInfo.Nodes = append(sw.hostInfo.Nodes[:i], sw.hostInfo.Nodes[i+1:]...)
			return
		}
	}
}

func (sw *SystemWatcher) executeTasks(tasks []string) {
	for _, task := range tasks {
		parts := strings.Split(task, " ")
		var err error
		res := &data.TaskResult{
			HostID: sw.hostInfo.HostID,
			Task:   task,
			Output: make([]byte, 0),
		}

		switch parts[0] {
		case utils.HOST_CMD_REBOOT:
			err = utils.RebootHost()

		case utils.HOST_CMD_UPDATE_OS:
			go func(task string) {
				err := utils.UpdateOS()
				utils.PostJsonHTTP(sw.appCfg.Server+utils.LISTEN_HOST_TASK_RESULT, &data.TaskResult{
					HostID: sw.hostInfo.HostID,
					Task:   task,
					Output: make([]byte, 0),
					Error:  err.Error(),
				})
			}(task)
			continue

		case utils.HOST_CMD_UPDATE_APP:
			go func(task string) {
				err := utils.SelfUpdate(utils.GITHUB_REPO)
				utils.PostJsonHTTP(sw.appCfg.Server+utils.LISTEN_HOST_TASK_RESULT, &data.TaskResult{
					HostID: sw.hostInfo.HostID,
					Task:   task,
					Output: make([]byte, 0),
					Error:  err.Error(),
				})
			}(task)
			continue

		case utils.HOST_CMD_EXEC:
			go func(task string) {
				output := make([]byte, 0)
				err := errors.New("no command specified")
				if len(parts) > 1 {
					output, err = exec.CommandContext(context.Background(), parts[1], parts[2:]...).Output()
				}
				utils.PostJsonHTTP(sw.appCfg.Server+utils.LISTEN_HOST_TASK_RESULT, &data.TaskResult{
					HostID: sw.hostInfo.HostID,
					Task:   task,
					Output: output,
					Error:  err.Error(),
				})
			}(task)
			continue
		}

		if err != nil {
			res.Error = err.Error()
		}
		utils.PostJsonHTTP(sw.appCfg.Server+utils.LISTEN_HOST_TASK_RESULT, res)
	}
}
