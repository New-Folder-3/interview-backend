package cmd

import (
	log "github.com/sirupsen/logrus"
	"interview-backend/internal/bootstrap"
	"interview-backend/internal/db"
	"interview-backend/util"
	"os"
	"path/filepath"
	"strconv"
)

var pid = -1
var pidFile string

func Init() {
	bootstrap.InitConfig()
	bootstrap.InitDB()
}

func Release() {
	db.Close()
}

func initDaemon() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	_ = os.MkdirAll(filepath.Join(exPath, "daemon"), 0700)
	pidFile = filepath.Join(exPath, "daemon/pid")
	if util.FileExist(pidFile) {
		bytes, err := os.ReadFile(pidFile)
		if err != nil {
			log.Fatal("failed to read pid file", err)
		}
		id, err := strconv.Atoi(string(bytes))
		if err != nil {
			log.Fatal("failed to parse pid data", err)
		}
		pid = id
	}
}
