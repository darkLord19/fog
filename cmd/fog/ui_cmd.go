package main

import (
	"fmt"
	"os"
	"time"

	"github.com/darkLord19/wtx/internal/daemon"
	fogenv "github.com/darkLord19/wtx/internal/env"
	"github.com/spf13/cobra"
)

var (
	uiPortFlag    int
	uiNoOpenFlag  bool
	uiTimeoutFlag time.Duration
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Open Fog web UI (starts fogd if needed)",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	uiCmd.Flags().IntVar(&uiPortFlag, "port", 8080, "fogd port")
	uiCmd.Flags().BoolVar(&uiNoOpenFlag, "no-open", false, "Do not auto-open browser")
	uiCmd.Flags().DurationVar(&uiTimeoutFlag, "timeout", 15*time.Second, "Wait timeout for fogd health")
	rootCmd.AddCommand(uiCmd)
}

func runUI() error {
	fogHome, err := fogenv.FogHome()
	if err != nil {
		return err
	}

	baseURL, err := daemon.EnsureRunning(fogHome, uiPortFlag, uiTimeoutFlag)
	if err != nil {
		return err
	}

	fmt.Printf("Fog UI: %s\n", baseURL)
	if uiNoOpenFlag {
		return nil
	}

	if err := daemon.OpenBrowser(baseURL); err != nil {
		fmt.Printf("Could not open browser automatically: %v\n", err)
		fmt.Printf("Open manually: %s\n", baseURL)
	}

	return nil
}
