package bootstrap

import (
	"interview-backend/cmd/flags"
	"interview-backend/util"
	"path/filepath"
)

func InitFile() {
	util.CreateFolder(filepath.Join(flags.DataDir, "image"))
	util.CreateFolder(filepath.Join(flags.DataDir, "audio"))
}
