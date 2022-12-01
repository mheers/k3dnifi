package cmd

import (
	"github.com/mheers/k3dnifi/helpers"
	"github.com/spf13/cobra"
)

var (
	restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "restart restarts a nifi cluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			nifiPod, err := getNifiPod()
			if err != nil {
				return err
			}
			return nifiPod.Restart()
		},
	}
)

func init() {
	registerNifiSelectorFlags(restartCmd)
}
