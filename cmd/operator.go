package cmd

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mheers/k3dnifi/helpers"
	"github.com/mheers/k3dnifi/nifi"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	operatorCmd = &cobra.Command{
		Use:   "operator",
		Short: "operator manages the nifi operator",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return cmd.Help()
		},
	}

	operatorVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "version reads the version of the nifi operator",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			operatorVersion, err := nifi.GetOperatorVersion()
			if err != nil {
				return err
			}
			cmd.Println(operatorVersion)
			return nil
		},
	}

	operatorInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "info reads the version and running infos of the nifi operator",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return renderOperator()
		},
	}

	operatorLogsCmd = &cobra.Command{
		Use:   "logs",
		Short: "logs reads the logs of the nifi operator",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			logs, err := nifi.GetOperatorLogs()
			if err != nil {
				return err
			}
			cmd.Println(logs)
			return nil
		},
	}

	operatorRestartCmd = &cobra.Command{
		Use:   "restart",
		Short: "restart restarts the nifi operator",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return nifi.RestartOperator()
		},
	}
)

func init() {
	operatorCmd.AddCommand(operatorLogsCmd)
	operatorCmd.AddCommand(operatorVersionCmd)
	operatorCmd.AddCommand(operatorInfoCmd)
	operatorCmd.AddCommand(operatorRestartCmd)
}

type operatorSummary struct {
	Name      string          `json:"name"`
	Namespace string          `json:"namespace"`
	Version   string          `json:"version"`
	Ready     string          `json:"ready"`
	Status    corev1.PodPhase `json:"status"`
	Created   time.Time       `json:"created"`
}

func podToSummary(operatorPod nifi.OperatorPod) []operatorSummary {
	var operatorSummaries []operatorSummary
	operatorSummaries = append(operatorSummaries, operatorSummary{
		Name:      operatorPod.Pod.Name,
		Namespace: operatorPod.Namespace,
		Version:   operatorPod.Version(),
		Ready:     operatorPod.Ready(),
		Status:    operatorPod.Status(),
		Created:   operatorPod.Created(),
	})
	return operatorSummaries
}

func renderOperator() error {
	operatorPod, err := nifi.GetOperatorPod()
	if err != nil {
		return err
	}
	summary := podToSummary(*operatorPod)
	if OutputFormatFlag == "table" {
		renderOperatorListTable(*operatorPod)
	}
	if OutputFormatFlag == "json" {
		err := helpers.PrintJSON(summary)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "yaml" {
		err := helpers.PrintYAML(summary)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "csv" {
		err := helpers.PrintCSV(summary)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderOperatorListTable(operatorPod nifi.OperatorPod) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Name", "Namespace", "Version", "Ready", "Status", "Created"})
	t.AppendRow(table.Row{operatorPod.Pod.Name, operatorPod.Namespace, operatorPod.Version(), operatorPod.Ready(), operatorPod.Status(), operatorPod.Created()})
	t.AppendSeparator()
	t.Render()
}
