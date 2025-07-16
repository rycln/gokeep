package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFilePath    = "urls"
	testDatabaseDsn = "test_dsn"
	testTimeout     = time.Duration(3) * time.Minute
	testKey         = "secret_key"
	testLoggerLevel = "info"
	testCfgFileName = "testcfg.json"
	testGRPCPort    = ":50052"
)

var testCfg = &Cfg{
	StorageFilePath: testFilePath,
	DatabaseDsn:     testDatabaseDsn,
	Timeout:         testTimeout,
	Key:             testKey,
	LogLevel:        testLoggerLevel,
	GRPCPort:        testGRPCPort,
}

func TestConfigBuilder_WithEnvParsing(t *testing.T) {
	t.Setenv("FILE_STORAGE_PATH", testCfg.StorageFilePath)
	t.Setenv("DATABASE_DSN", testCfg.DatabaseDsn)
	t.Setenv("TIMEOUT_DUR", testCfg.Timeout.String())
	t.Setenv("JWT_KEY", testCfg.Key)
	t.Setenv("LOG_LEVEL", testCfg.LogLevel)
	t.Setenv("GRPC_PORT", testGRPCPort)

	t.Run("valid test", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithEnvParsing().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}

func TestConfigBuilder_WithDefaultJWTKey(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithDefaultJWTKey().
			Build()
		assert.NoError(t, err)
		assert.NotEmpty(t, cfg.Key)
	})
}

func TestConfigBuilder_WithFlagParsing(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	t.Run("valid test", func(t *testing.T) {
		os.Args = []string{
			"./server",
			"-f=" + testCfg.StorageFilePath,
			"-d=" + testCfg.DatabaseDsn,
			"-o=" + testCfg.Timeout.String(),
			"-k=" + testCfg.Key,
			"-l=" + testCfg.LogLevel,
			"-g=" + testGRPCPort,
		}

		cfg, err := NewConfigBuilder().
			WithFlagParsing().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}

func TestConfigBuilder_WithConfigFile(t *testing.T) {
	file, err := os.Create(testCfgFileName)
	require.NoError(t, err)
	defer func() {
		err = file.Close()
		require.NoError(t, err)
		err = os.Remove(testCfgFileName)
		require.NoError(t, err)
	}()

	enc := json.NewEncoder(file)

	err = enc.Encode(&testCfg)
	require.NoError(t, err)

	t.Run("file name from flag", func(t *testing.T) {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()

		os.Args = []string{
			"./server",
			"-f=" + testCfg.StorageFilePath,
			"-d=" + testCfg.DatabaseDsn,
			"-t=" + testCfg.Timeout.String(),
			"-k=" + testCfg.Key,
			"-l=" + testCfg.LogLevel,
			"-c=" + testCfgFileName,
		}

		cfg, err := NewConfigBuilder().
			WithConfigFile().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})

	t.Run("file name from env", func(t *testing.T) {
		t.Setenv("CONFIG", testCfgFileName)

		cfg, err := NewConfigBuilder().
			WithConfigFile().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}
