package config

import (
	"encoding/json"
	"os"

	"github.com/stakingagency/nodemon/data"
)

func LoadNodeMonConfig(configFile string) (*data.NodeMonAppConfig, error) {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	cfg := &data.NodeMonAppConfig{}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

