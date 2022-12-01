package nifi

import (
	"fmt"
	"time"

	k3droot "github.com/mheers/k3droot/helpers"
	v1 "k8s.io/api/core/v1"
)

const (
	nifiKopImageSearchString = "konpyutaika/docker-images/nifikop"
	nifiKopContainerName     = "nifikop"
)

func GetOperatorPod() (*OperatorPod, error) {
	candidates, err := K8s.GetPodsByImage(nifiKopImageSearchString, false)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no nifikop pod found")
	}
	operatorPod := &OperatorPod{
		Pod:       candidates[0],
		Namespace: candidates[0].Namespace,
	}
	return operatorPod, nil
}

func ExecInOperator(args []string) error {
	operatorPod, err := GetOperatorPod()
	if err != nil {
		return err
	}
	podName := operatorPod.Pod.Name
	namespace := operatorPod.Namespace
	containerName := nifiKopContainerName
	return k3droot.ExecInNamespacePodContainer(namespace, podName, containerName, args)
}

func GetOperatorVersion() (string, error) {
	operatorPod, err := GetOperatorPod()
	if err != nil {
		return "", err
	}
	return operatorPod.Pod.Spec.Containers[0].Image, nil
}

func GetOperatorLogs() (string, error) {
	operatorPod, err := GetOperatorPod()
	if err != nil {
		return "", err
	}
	// Get the logs
	logs, err := K8s.GetLogsOfPod(*operatorPod.Pod)
	if err != nil {
		return "", err
	}
	return logs, nil
}

func RestartOperator() error {
	operatorPod, err := GetOperatorPod()
	if err != nil {
		return err
	}
	return K8s.DeletePod(*operatorPod.Pod)
}

type OperatorPod struct {
	Pod       *v1.Pod
	Namespace string
}

func (op *OperatorPod) FQN() string {
	return fmt.Sprintf("%s/%s", op.Namespace, op.Pod.Name)
}

func (op *OperatorPod) Ready() string {
	nReady := 0
	for _, containerStatus := range op.Pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			nReady++
		}
	}
	return fmt.Sprintf("%d/%d", nReady, len(op.Pod.Spec.Containers))
}

func (op *OperatorPod) Status() v1.PodPhase {
	return op.Pod.Status.Phase
}

func (op *OperatorPod) Created() time.Time {
	return op.Pod.CreationTimestamp.UTC()
}

func (op *OperatorPod) Version() string {
	operatorContainer, err := op.OperatorContainer()
	if err != nil {
		return ""
	}
	return operatorContainer.Image
}

func (op *OperatorPod) OperatorContainer() (*v1.Container, error) {
	for _, container := range op.Pod.Spec.Containers {
		if container.Name == "nifikop" {
			return &container, nil
		}
	}
	return nil, fmt.Errorf("no nifi container found in pod %s", op.Pod.Name)
}
