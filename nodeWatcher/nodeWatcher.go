package nodeWatcher

import (
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/pelletier/go-toml"
	"github.com/stakingagency/nodemon/data"
	"github.com/stakingagency/nodemon/systemWatcher"
)

type NodeWatcher struct {
	appCfg          *data.NodeMonAppConfig
	path            string
	listenInterface string
	sysWatcher      *systemWatcher.SystemWatcher
	prefs           *data.Preferences

	terminate bool
}

var log = logger.GetOrCreate("node-watcher")

func NewNodeWatcher(path string, listenInterface string, sysWatcher *systemWatcher.SystemWatcher) (*NodeWatcher, error) {
	n := &NodeWatcher{
		appCfg:          sysWatcher.GetAppConfig(),
		path:            path,
		listenInterface: listenInterface,
		sysWatcher:      sysWatcher,
		prefs:           &data.Preferences{},
	}

	return n, n.readNodePrefs()
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

func (nw *NodeWatcher) readNodePrefs() error {
	tree, err := toml.LoadFile(nw.path + "/config/prefs.toml")
	if err != nil {
		return err
	}

	return tree.Unmarshal(nw.prefs)
}
