package cmd

import (
	"errors"
	"fmt"

	"github.com/mheers/k3dnifi/helpers"
	"github.com/mheers/k3dnifi/nifi"
	k3droot "github.com/mheers/k3droot/helpers"
	"github.com/spf13/cobra"
)

var (
	clusterNamespace string
	clusterName      string
	nodeID           int

	shellCmd = &cobra.Command{
		Use:   "shell",
		Short: "shell into a nifi clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			nifiPod, err := getNifiPod()
			if err != nil {
				return err
			}
			podName := nifiPod.Pod.Name
			namespace := nifiPod.Pod.Namespace
			shell := defaultShell
			container, err := nifiPod.NifiContainer()
			if err != nil {
				return err
			}
			return k3droot.RootIntoNamespacePodContainer(namespace, podName, container.Name, shell)
		},
	}
)

func registerNifiSelectorFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&clusterNamespace, "namespace", "n", "", "namespace to use")
	cmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster to use")
	cmd.Flags().IntVarP(&nodeID, "node", "i", 0, "nifi cluster node to use")
}

func init() {
	registerNifiSelectorFlags(shellCmd)
}

func getNifiPod() (*nifi.NifiPod, error) {

	nifiPods, err := nifi.GetNifiPods()
	if err != nil {
		return nil, err
	}

	if len(nifiPods) == 0 {
		return nil, errors.New("no nifi pods found")
	}

	var nifiPod *nifi.NifiPod

	clusters := nifi.GetNifiClusters(nifiPods)
	namespaces := nifi.GetNifiNamespaces(nifiPods)

	if clusterName == "" {
		if len(clusters) == 1 {
			clusterName = clusters[0]
		} else {
			return nil, errors.New("multiple nifi clusters found, please specify one")
		}
	}

	if clusterNamespace == "" {
		if len(namespaces) == 1 {
			clusterNamespace = namespaces[0]
		} else {
			return nil, errors.New("multiple nifi namespaces found, please specify one")
		}
	}

	if nodeID == 0 {
		if len(nifiPods) == 1 {
			nodeID = nifiPods[0].NodeID()
		} else {
			return nil, errors.New("multiple nifi nodes found, please specify one")
		}
	}

	found := false
	for _, pod := range nifiPods {
		if pod.Namespace == clusterNamespace && pod.ClusterName == clusterName && pod.NodeID() == nodeID {
			nifiPod = pod
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("nifi cluster %s/%s not found", clusterNamespace, clusterName)
	}

	if nifiPod == nil {
		return nil, errors.New("no nifi cluster selected")
	}

	return nifiPod, nil
}
