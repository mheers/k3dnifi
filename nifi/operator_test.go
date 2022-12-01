package nifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOperatorPod(t *testing.T) {
	_, err := Init()
	require.NoError(t, err)

	operatorPod, err := GetOperatorPod()
	assert.Nil(t, err)
	assert.NotNil(t, operatorPod)
}
