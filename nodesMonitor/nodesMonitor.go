package nodesMonitor

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stakingagency/nodemon/data"
	"github.com/stakingagency/nodemon/nodeWatcher"
	"github.com/stakingagency/nodemon/systemWatcher"
)

type NodesMonitor struct {
	appCfg          *data.NodeMonAppConfig
	sysWatcher      *systemWatcher.SystemWatcher
	nodeWatchers    map[uint16]*nodeWatcher.NodeWatcher
	nodeWatchersMut sync.Mutex
}

var log = logger.GetOrCreate("nodes-monitor")

func NewNodesMonitor(appCfg *data.NodeMonAppConfig) (*NodesMonitor, error) {
	sysWatcher, err := systemWatcher.NewSystemWatcher(appCfg)
	if err != nil {
		return nil, err
	}

	m := &NodesMonitor{
		appCfg:       appCfg,
		sysWatcher:   sysWatcher,
		nodeWatchers: make(map[uint16]*nodeWatcher.NodeWatcher),
	}

	return m, nil
}

func (nm *NodesMonitor) StartTasks() {
	go nm.discoverNodes()
}

func (nm *NodesMonitor) discoverNodes() {
	firstTime := true
	for {
		time.Sleep(time.Second)
		nodeServices, err := getNodeServices()
		if err != nil {
			log.Error("get node services", "error", err)
			continue
		}

		nodeInterfaces := make(map[string]string)
		for nodePath, nodeService := range nodeServices {
			params := strings.Split(nodeService, " ")
			listenInterface := getNodeParam(params, "-rest-api-interface")
			if !strings.Contains(listenInterface, ":") {
				listenInterface = listenInterface + ":8080"
			}
			host := strings.Split(listenInterface, ":")[0]
			if host == "" || host == "0.0.0.0" {
				host = "localhost"
			}
			port := strings.Split(listenInterface, ":")[1]
			iPort, err := strconv.ParseUint(port, 10, 32)
			if err != nil {
				log.Warn("parse node port", "error", err, "line", nodeService)
				continue
			}

			found := false
			for _, oldInterface := range nodeInterfaces {
				oldPort, _ := strconv.ParseUint(strings.Split(oldInterface, ":")[1], 10, 32)
				if oldPort == iPort {
					log.Warn("duplicate port", "port", oldPort)
					found = true
					break
				}
			}
			if !found {
				nodeInterfaces[nodePath] = fmt.Sprintf("%s:%v", host, iPort)
			}
		}

		oldWatchers := nm.getNodeWatchers()
		for nodePath, newInterface := range nodeInterfaces {
			newPort, _ := strconv.ParseUint(strings.Split(newInterface, ":")[1], 10, 32)
			_, exists := oldWatchers[uint16(newPort)]
			if !exists {
				watcher, err := nodeWatcher.NewNodeWatcher(nodePath, newInterface, nm.sysWatcher)
				if err != nil {
					log.Warn("new node watcher", "error", err, "interface", newInterface)
					continue
				}

				nm.nodeWatchersMut.Lock()
				nm.nodeWatchers[uint16(newPort)] = watcher
				nm.nodeWatchersMut.Unlock()

				nm.sysWatcher.AddNode(newInterface)
				log.Info("new node", "interface", newInterface)

				watcher.StartTasks()
			}
		}

		for oldPort, oldWatcher := range oldWatchers {
			found := false
			for _, newInterface := range nodeInterfaces {
				newPort, _ := strconv.ParseUint(strings.Split(newInterface, ":")[1], 10, 32)
				if oldPort == uint16(newPort) {
					found = true
					break
				}
			}
			if !found {
				oldWatcher.StopTasks()
				nm.sysWatcher.RemoveNode(oldWatcher.GetListenInterface())
				log.Warn("node removed", "interface", oldWatcher.GetListenInterface())

				nm.nodeWatchersMut.Lock()
				delete(nm.nodeWatchers, oldPort)
				nm.nodeWatchersMut.Unlock()
			}
		}

		if firstTime {
			firstTime = false
			nm.sysWatcher.StartTasks()
		}
		time.Sleep(time.Minute)
	}
}

func (nm *NodesMonitor) getNodeWatchers() map[uint16]*nodeWatcher.NodeWatcher {
	res := make(map[uint16]*nodeWatcher.NodeWatcher)
	nm.nodeWatchersMut.Lock()
	for port, watcher := range nm.nodeWatchers {
		res[port] = watcher
	}
	nm.nodeWatchersMut.Unlock()

	return res
}

func getNodeServices() (map[string]string, error) {
	res := make(map[string]string)
	rootPath := "/etc/systemd/system/"
	dirEntries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		path := rootPath + entry.Name()
		if !strings.HasSuffix(path, ".service") || entry.IsDir() {
			continue
		}

		bytes, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		lines := strings.Split(strings.ReplaceAll(string(bytes), "\r", ""), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "ExecStart=") {
				line = strings.TrimPrefix(line, "ExecStart=")
				appPath := strings.Split(line, " ")[0]
				appPathParts := strings.Split(appPath, "/")
				app := appPathParts[len(appPathParts)-1]
				if app == "node" {
					res[strings.TrimSuffix(appPath, "/node")] = strings.TrimSpace(strings.TrimPrefix(line, appPath))
				}
			}
		}
	}

	return res, err
}

func getNodeParam(params []string, param string) string {
	n := len(params)
	for i := 0; i < n; i++ {
		if strings.HasSuffix(params[i], param) && i < n-1 {
			return params[i+1]
		}
	}

	return ""
}
