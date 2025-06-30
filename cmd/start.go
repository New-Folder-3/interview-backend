package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Silent start",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	initDaemon()
	if pid != -1 {
		_, err := os.FindProcess(pid)
		if err == nil {
			log.Info("alist already started, pid ", pid)
			return
		}
	}
	args := os.Args
	args[1] = "server"
	cmd := &exec.Cmd{
		Path: args[0],
		Args: args,
		Env:  os.Environ(),
	}
	stdout, err := os.OpenFile(filepath.Join(filepath.Dir(pidFile), "start.log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(os.Getpid(), ": failed to open start log file:", err)
	}
	cmd.Stderr = stdout
	cmd.Stdout = stdout
	err = cmd.Start()
	if err != nil {
		log.Fatal("failed to start children process: ", err)
	}
	log.Infof("success start pid: %d", cmd.Process.Pid)
	err = os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0666)
	if err != nil {
		log.Warn("failed to record pid, you may not be able to stop the program with command")
	}
}

func init() {
	RootCmd.AddCommand(StartCmd)
}
