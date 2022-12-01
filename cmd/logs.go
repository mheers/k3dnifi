package cmd

import (
	"github.com/mheers/k3dnifi/helpers"
	"github.com/spf13/cobra"
)

var (
	logsCmd = &cobra.Command{
		Use:   "logs",
		Short: "logs manages the logs of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return cmd.Help()
		},
	}

	logsLsCmd = &cobra.Command{
		Use:     "ls",
		Short:   "ls lists the logs of nifi clusters",
		Aliases: []string{"list"},
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return listLogs()
		},
	}

	logsAppCmd = &cobra.Command{
		Use:   "app",
		Short: "app reads the app logs of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getLogs("app")
		},
	}

	logsBootstrapCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap reads the bootstrap logs of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getLogs("bootstrap")
		},
	}

	logsDeprecationCmd = &cobra.Command{
		Use:   "deprecation",
		Short: "deprecation reads the deprecation logs of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getLogs("deprecation")
		},
	}

	logsRequestCmd = &cobra.Command{
		Use:   "request",
		Short: "request reads the request logs of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getLogs("request")
		},
	}

	logsUserCmd = &cobra.Command{
		Use:   "user",
		Short: "user reads the user logs of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getLogs("user")
		},
	}
)

func init() {

	logsCmd.AddCommand(logsLsCmd)
	registerNifiSelectorFlags(logsLsCmd)

	logsCmd.AddCommand(logsAppCmd)
	registerNifiSelectorFlags(logsAppCmd)

	logsCmd.AddCommand(logsBootstrapCmd)
	registerNifiSelectorFlags(logsBootstrapCmd)

	logsCmd.AddCommand(logsDeprecationCmd)
	registerNifiSelectorFlags(logsDeprecationCmd)

	logsCmd.AddCommand(logsRequestCmd)
	registerNifiSelectorFlags(logsRequestCmd)

	logsCmd.AddCommand(logsUserCmd)
	registerNifiSelectorFlags(logsUserCmd)
}

func getLogs(logName string) error {
	nifiPod, err := getNifiPod()
	if err != nil {
		return err
	}
	return nifiPod.GetLogs(logName)
}

func listLogs() error {
	nifiPod, err := getNifiPod()
	if err != nil {
		return err
	}
	return nifiPod.ListLogs()
}
