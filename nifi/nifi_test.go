package nifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNifiPods(t *testing.T) {
	_, err := Init()
	require.NoError(t, err)

	nifiPods, err := GetNifiPods()
	assert.Nil(t, err)
	assert.NotNil(t, nifiPods)
}

func TestLogsPVHostPath(t *testing.T) {
	_, err := Init()
	require.NoError(t, err)

	nifiPods, err := GetNifiPods()
	require.NoError(t, err)
	assert.NotNil(t, nifiPods)

	pod := nifiPods[0]
	path, err := pod.LogsPVHostPath()
	require.NoError(t, err)
	assert.NotEmpty(t, path)
}
