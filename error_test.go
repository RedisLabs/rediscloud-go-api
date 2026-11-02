package rediscloud_api

import (
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertInternalError(t *testing.T, err error, expected *internal.Error) {
	t.Helper()

	var actual *internal.Error
	require.ErrorAs(t, err, &actual)
	assert.Equal(t, expected, actual)
}
