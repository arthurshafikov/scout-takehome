package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomString(t *testing.T) {
	result, err := GenerateRandomString(20)

	require.NoError(t, err)
	require.Len(t, result, 20)
}
