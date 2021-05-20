package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Env        string `toml:"env"`
	BindAddr   string `toml:"apiserver_port"`
	SendToChat bool   `toml:"send_to_chat"`
	BasePath   string
}

var (
	Settings *Config

	defaultWorkingDir, _ = os.Getwd()
	defaultconfigPath    = "configs/config.toml"
)

func init() {
	workingDir := os.Getenv("WORKING_DIR")
	if workingDir == "" {
		workingDir = defaultWorkingDir
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultconfigPath
	}

	Settings = new(Config)
	_, err := toml.DecodeFile(filepath.Join(workingDir, configPath), Settings)

	if err != nil {
		panic(fmt.Sprintf("error with decoding config file: %s", err))
	}

	Settings.BasePath = workingDir
}
