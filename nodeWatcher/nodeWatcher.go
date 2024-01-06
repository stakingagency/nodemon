package nodeWatcher

import (
	"os"
	"strings"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stakingagency/nodemon/data"
	"github.com/stakingagency/nodemon/systemWatcher"
)

type NodeWatcher struct {
	appCfg          *data.NodeMonAppConfig
	path            string
	listenInterface string
	sysWatcher      *systemWatcher.SystemWatcher

	identity  string
	terminate bool
}

var log = logger.GetOrCreate("node-watcher")

func NewNodeWatcher(path string, listenInterface string, sysWatcher *systemWatcher.SystemWatcher) (*NodeWatcher, error) {
	n := &NodeWatcher{
		appCfg:          sysWatcher.GetAppConfig(),
		path:            path,
		listenInterface: listenInterface,
		sysWatcher:      sysWatcher,
	}

	n.readIdentity()

	return n, nil
}

func (nw *NodeWatcher) StartTasks() {
	go nw.watchNode()
}

func (nw *NodeWatcher) StopTasks() {
	nw.terminate = true
}

func (nw *NodeWatcher) GetListenInterface() string {
	return nw.listenInterface
}

func (nw *NodeWatcher) readIdentity() error {
	bytes, err := os.ReadFile(nw.path + "/config/prefs.toml")
	if err != nil {
		return err
	}

	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(strings.ToLower(line), "identity") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				nw.identity = strings.TrimSpace(strings.ReplaceAll(parts[1], "\"", ""))
				break
			}
		}
	}

	return nil
}
