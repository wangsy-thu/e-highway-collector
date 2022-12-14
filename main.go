package main

import (
	"e-highway-collector/config"
	"e-highway-collector/lib/logger"
	"e-highway-collector/tcp"
	"fmt"
	"os"
)

const configFile string = "highway.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 13000,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "highway-collector",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		},
		tcp.MakeEchoHandler())

	if err != nil {
		logger.Error(err)
	}
}
