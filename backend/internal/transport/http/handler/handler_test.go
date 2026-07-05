package handler

import (
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var (
	key      = "someParam"
	expected = 22
)

func TestParseIntegerFromParam(t *testing.T) {
	handler := &Handler{}

	result, err := handler.parseIntegerFromParam(&gin.Context{
		Params: gin.Params{
			{
				Key:   key,
				Value: "22",
			},
		},
	}, key)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestParseIntegerFromParamNotFoundParam(t *testing.T) {
	handler := &Handler{}

	result, err := handler.parseIntegerFromParam(&gin.Context{}, key)

	require.Equal(t, err.Error(), "the param someParam is missing")
	require.Zero(t, result)
}

func TestParseIntegerFromParamCouldntConvert(t *testing.T) {
	handler := &Handler{}

	result, err := handler.parseIntegerFromParam(&gin.Context{
		Params: gin.Params{
			{
				Key:   key,
				Value: "WrongParam",
			},
		},
	}, key)

	require.ErrorIs(t, err, strconv.ErrSyntax)
	require.Zero(t, result)
}
