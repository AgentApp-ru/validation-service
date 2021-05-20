package config

import (
    "flag"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Env        string `toml:"env"`
    BindAddr   string `toml:"apiserver_port"`
    SendToChat bool   `toml:"send_to_chat"`
    BasePath   string `toml:"base_path"`
}

var (
    Settings   *Config
    configPath string
)

func init() {
    flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
    flag.Parse()

    Settings = new(Config)
    _, err := toml.DecodeFile(configPath, Settings)
    if err != nil {
        panic(err)
    }
}
