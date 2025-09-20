package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
	require.Equal(t, 1, CompareVersions("1.0.1", "1.0.0"))
	require.Equal(t, 0, CompareVersions("v1.0.0", "v1.0.0"))
}

func TestCanUseGoVersion(t *testing.T) {
	require.True(t, CanUseGoVersion("1.20", "1.21"))
	require.False(t, CanUseGoVersion("1.21", "1.20"))
}
