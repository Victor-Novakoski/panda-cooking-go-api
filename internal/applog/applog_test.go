package applog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"panda-cooking-go-api/internal/applog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerAddsRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := applog.New(true, &buf).With("app", "api")

	ctx := applog.WithRequestID(context.Background(), "abc123")
	log.InfoContext(ctx, "teste")

	var line map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &line))
	assert.Equal(t, "abc123", line["request_id"])
	assert.Equal(t, "api", line["app"])
	assert.Equal(t, "teste", line["msg"])
}
