package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"interview-backend/cmd/flags"
	"os"
	"path/filepath"
)

var RootCmd = &cobra.Command{
	Use:   "interview",
	Short: "interview",
	Long:  "The backend of an AI-powered interviewee-assistant software",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)
	RootCmd.PersistentFlags().StringVar(&flags.DataDir, "data", exeDir+"/data", "data folder")
	RootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "start with debug mode")
	RootCmd.PersistentFlags().BoolVar(&flags.Dev, "dev", false, "start with dev mode")
}
