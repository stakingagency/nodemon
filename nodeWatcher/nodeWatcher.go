package nodeWatcher

import (
	"os"
	"strings"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/pelletier/go-toml"
	"github.com/stakingagency/nodemon/data"
	"github.com/stakingagency/nodemon/systemWatcher"
)

type NodeWatcher struct {
	appCfg           *data.NodeMonAppConfig
	path             string
	listenInterface  string
	sysWatcher       *systemWatcher.SystemWatcher
	prefs            *data.Preferences
	validatorKey     string
	allValidatorKeys []string
	isMultiKey       bool

	terminate bool
}

var log = logger.GetOrCreate("node-watcher")

func NewNodeWatcher(path string, listenInterface string, sysWatcher *systemWatcher.SystemWatcher) (*NodeWatcher, error) {
	n := &NodeWatcher{
		appCfg:           sysWatcher.GetAppConfig(),
		path:             path,
		listenInterface:  listenInterface,
		sysWatcher:       sysWatcher,
		prefs:            &data.Preferences{},
		allValidatorKeys: make([]string, 0),
	}

	err := n.readNodePrefs()
	if err != nil {
		return nil, err
	}

	err = n.readValidatorKeys()
	if err != nil {
		return nil, err
	}

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

func (nw *NodeWatcher) readNodePrefs() error {
	tree, err := toml.LoadFile(nw.path + "/config/prefs.toml")
	if err != nil {
		return err
	}

	return tree.Unmarshal(nw.prefs)
}

func (nw *NodeWatcher) readValidatorKeys() error {
	bytes, err := os.ReadFile(nw.path + "/config/validatorKey.pem")
	if err != nil {
		return err
	}

	beginKey := "BEGIN PRIVATE KEY for "
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.ReplaceAll(strings.TrimSpace(line), "-----", "")
		if strings.HasPrefix(line, beginKey) {
			nw.validatorKey = strings.TrimPrefix(line, beginKey)
			break
		}
	}

	bytes, err = os.ReadFile(nw.path + "/config/allValidatorsKeys.pem")
	if err != nil {
		return nil
	}

	lines = strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.ReplaceAll(strings.TrimSpace(line), "-----", "")
		if strings.HasPrefix(line, beginKey) {
			nw.allValidatorKeys = append(nw.allValidatorKeys, strings.TrimPrefix(line, beginKey))
		}
	}

	nw.isMultiKey = len(nw.allValidatorKeys) > 0

	return nil
}
