package nifi

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mheers/k3dnifi/helpers"
	k3droot "github.com/mheers/k3droot/helpers"
	corev1 "k8s.io/api/core/v1"
)

const (
	volumeNameConf = "conf"
	volumeNameLogs = "logs"
	// volumeNameFlowfile   = "flowfile-repository"
	// volumeNameContent    = "content-repository"
	// volumeNameProvenance = "provenance-repository"
	// volumeNameStatus     = "status-repository"
	// volumeNameState      = "state"
	// volumeNameDatabase   = "database-repository"
	// volumeNameData       = "data"
)

type K8sClient struct {
	k3droot.K8sClient
}

var (
	K8s  = &K8sClient{}
	Ctx  = context.Background()
	once sync.Once
)

func Init() (*K8sClient, error) {
	var errR error
	once.Do(func() {
		k8sclient, err := k3droot.Init()
		K8s.K8sClient = *k8sclient
		errR = err
	})
	return K8s, errR
}

type NifiPod struct {
	Pod         *corev1.Pod
	ClusterName string
	Namespace   string
}

func (nfp *NifiPod) FQN() string {
	return fmt.Sprintf("%s/%s/%s", nfp.Namespace, nfp.ClusterName, nfp.Pod.Name)
}
func (nfp *NifiPod) Ready() string {
	nReady := 0
	for _, containerStatus := range nfp.Pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			nReady++
		}
	}
	return fmt.Sprintf("%d/%d", nReady, len(nfp.Pod.Spec.Containers))
}
func (nfp *NifiPod) Status() string {
	if nfp.Pod.ObjectMeta.DeletionTimestamp != nil {
		return "Terminating"
	}
	return string(nfp.Pod.Status.Phase)
}

func (nfp *NifiPod) Created() time.Time {
	return nfp.Pod.CreationTimestamp.UTC()
}

func (nfp *NifiPod) Age() time.Duration {
	return time.Now().UTC().Sub(nfp.Created())
}

func (nfp *NifiPod) NodeID() int {
	nodeID, err := strconv.Atoi(nfp.Pod.Labels["nodeId"])
	if err != nil {
		return -1
	}
	return nodeID
}

func (nfp *NifiPod) Version() string {
	nifiContainer, err := nfp.NifiContainer()
	if err != nil {
		return ""
	}
	return nifiContainer.Image
}

func (nfp *NifiPod) pvHostPath(volumeName string) (string, error) {
	for _, volumeMount := range nfp.Pod.Spec.Volumes {
		if volumeMount.Name == volumeName {
			hostPath, err := K8s.K8sClient.GetHostPathOfVolumeMount(nfp.Namespace, volumeMount)
			if err != nil {
				if err.Error() == "volumeName is empty" {
					return "", fmt.Errorf("no volume found for %s", volumeName)
				}
				return "", err
			}
			return hostPath, nil
		}
	}
	return "", errors.New("no hostpath found for volume " + volumeName)
}

func (nfp *NifiPod) LogsPVHostPath() (string, error) {
	volumeName := volumeNameLogs
	return nfp.pvHostPath(volumeName)
}

func (nfp *NifiPod) ConfPVHostPath() (string, error) {
	volumeName := volumeNameConf
	return nfp.pvHostPath(volumeName)
}

func (nfp *NifiPod) Exec(args []string) error {
	podName := nfp.Pod.Name
	namespace := nfp.Namespace
	container, err := nfp.NifiContainer()
	if err != nil {
		return err
	}
	return k3droot.ExecInNamespacePodContainer(namespace, podName, container.Name, args)
}

func (nfp *NifiPod) GetLogs(logName string) error {
	if nfp.Status() == "Running" {
		return nfp.Exec([]string{"cat", fmt.Sprintf("/opt/nifi/nifi-current/logs/nifi-%s.log", logName)})
	} else {
		path, err := nfp.LogsPVHostPath()
		if err != nil {
			return err
		}
		return k3droot.RunInNodeOfPod(*nfp.Pod, []string{"cat", fmt.Sprintf("%s/nifi-%s.log", path, logName)})
	}
}

func (nfp *NifiPod) ListLogs() error {
	if nfp.Status() == "Running" {
		return nfp.Exec([]string{"ls", "/opt/nifi/nifi-current/logs/"})
	} else {
		path, err := nfp.LogsPVHostPath()
		if err != nil {
			return err
		}
		return k3droot.RunInNodeOfPod(*nfp.Pod, []string{"ls", fmt.Sprintf("%s/", path)})
	}
}

func (nfp *NifiPod) GetConf(logName string) error {
	if nfp.Status() == "Running" {
		return nfp.Exec([]string{"cat", fmt.Sprintf("/opt/nifi/nifi-current/conf/%s", logName)})
	} else {
		path, err := nfp.ConfPVHostPath()
		if err != nil {
			return err
		}
		return k3droot.RunInNodeOfPod(*nfp.Pod, []string{"cat", fmt.Sprintf("%s/%s", path, logName)})
	}
}

func (nfp *NifiPod) ListConf() error {
	if nfp.Status() == "Running" {
		return nfp.Exec([]string{"ls", "/opt/nifi/nifi-current/conf/"})
	} else {
		path, err := nfp.ConfPVHostPath()
		if err != nil {
			return err
		}
		return k3droot.RunInNodeOfPod(*nfp.Pod, []string{"ls", fmt.Sprintf("%s/", path)})
	}
}

func (nfp *NifiPod) NifiContainer() (*corev1.Container, error) {
	for _, container := range nfp.Pod.Spec.Containers {
		if container.Name == "nifi" {
			return &container, nil
		}
	}
	return nil, fmt.Errorf("no nifi container found in pod %s", nfp.Pod.Name)
}

func (nfp *NifiPod) Restart() error {
	return K8s.DeletePod(*nfp.Pod)
}

func GetNifiPods() ([]*NifiPod, error) {
	nifiPods := []*NifiPod{}

	pods, err := K8s.K8sClient.GetAllPods()
	if err != nil {
		return nil, err
	}

	for _, pod := range pods {
		// check if OwnerReferences is kind NifiCluster
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "NifiCluster" {
				nifiPods = append(nifiPods, &NifiPod{
					Pod:         pod,
					ClusterName: owner.Name,
					Namespace:   pod.Namespace,
				})
			}
		}
	}
	return nifiPods, nil
}

func GetNifiClusters(nifiPods []*NifiPod) []string {
	clusterNames := []string{}
	for _, nifiPod := range nifiPods {
		if _, found := helpers.Find(clusterNames, nifiPod.ClusterName); !found {
			clusterNames = append(clusterNames, nifiPod.ClusterName)
		}
	}
	return clusterNames
}

func GetNifiNamespaces(nifiPods []*NifiPod) []string {
	namespaces := []string{}
	for _, nifiPod := range nifiPods {
		if _, found := helpers.Find(namespaces, nifiPod.Namespace); !found {
			namespaces = append(namespaces, nifiPod.Namespace)
		}
	}
	return namespaces
}
