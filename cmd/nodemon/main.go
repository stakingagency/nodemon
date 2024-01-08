package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
	"github.com/stakingagency/nodemon/config"
	"github.com/stakingagency/nodemon/nodesMonitor"
	"github.com/stakingagency/nodemon/utils"
)

var (
	log        = logger.GetOrCreate("nodemon")
	appVersion = "v0.0.0"
)

func main() {
	err := startLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Info(os.Args[0], "version", appVersion)
	utils.AppVersion = appVersion

	appCfg, err := config.LoadNodeMonConfig("config.json")
	if err != nil {
		log.Error("load config", "error", err)
		os.Exit(1)
	}

	nodesMonitor, err := nodesMonitor.NewNodesMonitor(appCfg)
	if err != nil {
		log.Error("new nodes monitor", "error", err)
		os.Exit(1)
	}

	nodesMonitor.StartTasks()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sigs:
			log.Info("terminating at user's signal...")
			os.Exit(0)
		}
	}
}

func startLogger() error {
	err := logger.SetLogLevel("*:" + logger.LogDebug.String())
	if err != nil {
		return err
	}

	args := file.ArgsFileLogging{
		WorkingDir:      ".",
		DefaultLogsPath: "logs",
		LogFilePrefix:   "nodesmon",
	}
	fileLogging, err := file.NewFileLogging(args)
	if err != nil {
		return fmt.Errorf("%w creating a log file", err)
	}

	err = fileLogging.ChangeFileLifeSpan(time.Hour*24, 1024)
	if err != nil {
		return err
	}

	return nil
}
