package interceptors

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInterceptorLogger(t *testing.T) {
	t.Run("logs with different levels and field types", func(t *testing.T) {
		observedZapCore, observedLogs := observer.New(zap.DebugLevel)
		logger := zap.New(observedZapCore)

		interceptorLogger := InterceptorLogger(logger)

		testCases := []struct {
			name     string
			level    logging.Level
			msg      string
			fields   []any
			expected zapcore.Level
		}{
			{
				name:     "debug level",
				level:    logging.LevelDebug,
				msg:      "debug message",
				fields:   []any{"debug_key", "debug_value"},
				expected: zap.DebugLevel,
			},
			{
				name:     "info level",
				level:    logging.LevelInfo,
				msg:      "info message",
				fields:   []any{"info_key", 42},
				expected: zap.InfoLevel,
			},
			{
				name:     "warn level",
				level:    logging.LevelWarn,
				msg:      "warn message",
				fields:   []any{"warn_key", true},
				expected: zap.WarnLevel,
			},
			{
				name:     "error level",
				level:    logging.LevelError,
				msg:      "error message",
				fields:   []any{"error_key", 3.14},
				expected: zap.ErrorLevel,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				interceptorLogger.Log(context.Background(), tc.level, tc.msg, tc.fields...)

				logs := observedLogs.TakeAll()
				require.Len(t, logs, 1, "expected exactly one log entry")

				log := logs[0]
				assert.Equal(t, tc.msg, log.Message)
				assert.Equal(t, tc.expected, log.Level)

				require.Len(t, log.Context, 1, "expected exactly one field")
				assert.Equal(t, tc.fields[0].(string), log.Context[0].Key)
			})
		}
	})

	t.Run("panics on unknown level", func(t *testing.T) {
		observedZapCore, _ := observer.New(zap.DebugLevel)
		logger := zap.New(observedZapCore)
		interceptorLogger := InterceptorLogger(logger)

		assert.PanicsWithValue(t, "unknown level 99", func() {
			interceptorLogger.Log(context.Background(), 99, "test message")
		})
	})

	t.Run("handles empty fields", func(t *testing.T) {
		observedZapCore, observedLogs := observer.New(zap.DebugLevel)
		logger := zap.New(observedZapCore)
		interceptorLogger := InterceptorLogger(logger)

		interceptorLogger.Log(context.Background(), logging.LevelInfo, "no fields")

		logs := observedLogs.TakeAll()
		require.Len(t, logs, 1)
		assert.Empty(t, logs[0].Context)
	})
}
