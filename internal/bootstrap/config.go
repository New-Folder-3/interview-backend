package bootstrap

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"interview/cmd/flags"
	"interview/internal/conf"
	"interview/util"
	"os"
	"path/filepath"
)

func InitConfig() {
	configPath := filepath.Join(flags.DataDir, "config.json")
	conf.Conf = conf.DefaultConfig()
	if !util.FileExist(configPath) {
		log.Infof("config file not exists, setting default config file")
		_, err := util.CreateFile(configPath)
		if err != nil {
			log.Fatalf("init config file failed: %s", err)
		}
	} else {
		configBytes, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("load config file failed: %s", err)
		}
		if err := json.Unmarshal(configBytes, conf.Conf); err != nil {
			log.Fatalf("init config file failed: %s", err)
		}
	}
	if err := util.JsonToFile(configPath, conf.Conf); err != nil {
		log.Fatalf("write config file failed: %s", err)
	}
}
