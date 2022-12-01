package cmd

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mheers/k3dnifi/helpers"
	"github.com/mheers/k3dnifi/nifi"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "list all nifi clusters",
		Aliases: []string{"ls"},
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			nifiPods, err := nifi.GetNifiPods()
			if err != nil {
				panic(err)
			}

			return renderClusters(nifiPods)
		},
	}
)

type nifiClusterSummary struct {
	NifiNode    string        `json:"nifiNode"`
	NodeID      int           `json:"nodeId"`
	ClusterName string        `json:"clusterName"`
	Namespace   string        `json:"namespace"`
	Version     string        `json:"version"`
	Ready       string        `json:"ready"`
	Status      string        `json:"status"`
	Created     time.Time     `json:"created"`
	Age         time.Duration `json:"age"`
}

func podsToSummary(nifiPods []*nifi.NifiPod) []*nifiClusterSummary {
	var nifiClusterSummaries []*nifiClusterSummary
	for _, nifiPod := range nifiPods {
		nifiClusterSummaries = append(nifiClusterSummaries, &nifiClusterSummary{
			NifiNode:    nifiPod.Pod.Name,
			NodeID:      nifiPod.NodeID(),
			ClusterName: nifiPod.ClusterName,
			Namespace:   nifiPod.Namespace,
			Version:     nifiPod.Version(),
			Ready:       nifiPod.Ready(),
			Status:      nifiPod.Status(),
			Created:     nifiPod.Created(),
			Age:         nifiPod.Age(),
		})
	}
	return nifiClusterSummaries
}

func renderClusters(nifiPods []*nifi.NifiPod) error {
	summary := podsToSummary(nifiPods)
	if OutputFormatFlag == "table" {
		renderListTable(nifiPods)
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

func renderListTable(nifiPod []*nifi.NifiPod) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"NifiNode", "NodeID", "ClusterName", "Namespace", "Version", "Ready", "Status", "Created", "Age"})
	for _, nifiPod := range nifiPod {
		t.AppendRow(table.Row{nifiPod.Pod.Name, nifiPod.NodeID(), nifiPod.ClusterName, nifiPod.Namespace, nifiPod.Version(), nifiPod.Ready(), nifiPod.Status(), nifiPod.Created(), nifiPod.Age()})
		t.AppendSeparator()
	}
	t.Render()
}
