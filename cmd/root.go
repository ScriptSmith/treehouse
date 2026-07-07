package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

var rootCmd = &cobra.Command{
	Use:   "treehouse",
	Short: "Manage a pool of git worktrees for parallel AI agent workflows",
	Long: `Treehouse maintains a pool of reusable, pre-warmed git worktrees
so that multiple AI coding agents can work on the same repo in parallel.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getRunE(cmd, args)
	},
}

func init() {
	rootCmd.SetVersionTemplate(`{{.Version}}` + "\n")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
