package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// LogLevelFlag describes the verbosity of logs
	LogLevelFlag string
	// ConfigFileFlag holds the path to the config file
	ConfigFileFlag string

	// OutputFormatFlag can be json, yaml or table
	OutputFormatFlag string

	// // Config holds the read config
	// Config *config.Config

	defaultShell = "bash"

	rootCmd = &cobra.Command{
		Use:   "k3dnifi",
		Short: "k3dnifi is a command line interface for nifi in k3d",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&LogLevelFlag, "log-level", "l", "error", "possible values are debug, error, fatal, panic, info, trace")
	rootCmd.PersistentFlags().StringVarP(&OutputFormatFlag, "output-format", "o", "table", "format [json|table|yaml|csv]")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(operatorCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(confCmd)
	rootCmd.AddCommand(logsCmd)
}
