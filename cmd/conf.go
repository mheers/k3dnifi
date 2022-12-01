package cmd

import (
	"github.com/mheers/k3dnifi/helpers"
	"github.com/spf13/cobra"
)

var (
	confCmd = &cobra.Command{
		Use:   "conf",
		Short: "conf manages the configuration of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return cmd.Help()
		},
	}

	confLsCmd = &cobra.Command{
		Use:     "ls",
		Short:   "ls lists the configuration of nifi clusters",
		Aliases: []string{"list"},
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return listConf()
		},
	}

	confAuthorizationsCmd = &cobra.Command{
		Use:   "authorizations",
		Short: "authorizations reads the authorizations of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getConf("authorizations.xml")
		},
	}

	confAuthorizersCmd = &cobra.Command{
		Use:   "authorizers",
		Short: "authorizers reads the authorizers of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getConf("authorizers.xml")
		},
	}

	confBootstrapCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap reads the bootstrap of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getConf("bootstrap.conf")
		},
	}

	confLogbackCmd = &cobra.Command{
		Use:   "logback",
		Short: "logback reads the logback of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getConf("logback.xml")
		},
	}

	confUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "users reads the users of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getConf("users.xml")
		},
	}

	confPropertiesCmd = &cobra.Command{
		Use:   "properties",
		Short: "properties manages the properties of nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return getConf("nifi.properties")
		},
	}
)

func init() {
	confCmd.AddCommand(confPropertiesCmd)
	registerNifiSelectorFlags(confPropertiesCmd)

	confCmd.AddCommand(confLsCmd)
	registerNifiSelectorFlags(confLsCmd)

	confCmd.AddCommand(confAuthorizationsCmd)
	registerNifiSelectorFlags(confAuthorizationsCmd)

	confCmd.AddCommand(confAuthorizersCmd)
	registerNifiSelectorFlags(confAuthorizersCmd)

	confCmd.AddCommand(confBootstrapCmd)
	registerNifiSelectorFlags(confBootstrapCmd)
	confCmd.AddCommand(confLogbackCmd)

	registerNifiSelectorFlags(confLogbackCmd)

	confCmd.AddCommand(confUsersCmd)
	registerNifiSelectorFlags(confUsersCmd)
}

func getConf(confName string) error {
	nifiPod, err := getNifiPod()
	if err != nil {
		return err
	}
	return nifiPod.GetConf(confName)
}

func listConf() error {
	nifiPod, err := getNifiPod()
	if err != nil {
		return err
	}
	return nifiPod.ListConf()
}
