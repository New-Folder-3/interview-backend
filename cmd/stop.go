package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"syscall"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop interview server by daemon/pid file",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
	},
}

func stop() {
	initDaemon()
	if pid == -1 {
		log.Info("Seems not have been started. Try use `interview start` to start server.")
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Errorf("failed to find process by pid: %d, reason: %v", pid, process)
		return
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Errorf("failed to terminate process %d: %v", pid, err)
	} else {
		log.Info("terminated process: ", pid)
	}
	err = os.Remove(pidFile)
	if err != nil {
		log.Errorf("failed to remove pid file")
	}
	pid = -1
}

func init() {
	RootCmd.AddCommand(StopCmd)
}
