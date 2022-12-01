package cmd

import (
	"github.com/mheers/k3dnifi/helpers"
	"github.com/spf13/cobra"
)

var (
	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "exec into a nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return execInNifi(args)
		},
	}
)

func init() {
	registerNifiSelectorFlags(execCmd)
}

func execInNifi(args []string) error {
	nifiPod, err := getNifiPod()
	if err != nil {
		return err
	}
	return nifiPod.Exec(args)
}
