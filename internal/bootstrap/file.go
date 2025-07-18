package bootstrap

import (
	"interview/cmd/flags"
	"interview/util"
	"path/filepath"
)

func InitFile() {
	util.CreateFolder(filepath.Join(flags.DataDir, "image"))
	util.CreateFolder(filepath.Join(flags.DataDir, "audio"))
	util.CreateFolder(filepath.Join(flags.DataDir, "video"))
}
