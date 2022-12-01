package cmd

import (
	"testing"

	"github.com/mheers/k3dnifi/nifi"
	"github.com/stretchr/testify/require"
)

func TestGetNifiPod(t *testing.T) {
	_, err := nifi.Init()
	require.NoError(t, err)

	pod, err := getNifiPod()
	require.NoError(t, err)
	require.NotNil(t, pod)
}
