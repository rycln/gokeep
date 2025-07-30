package bin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadFile(t *testing.T) {
	testString := "test content"
	encodedContent := base64.StdEncoding.EncodeToString([]byte(testString))
	testJSON := `{"bin":"` + encodedContent + `"}`
	testContent := []byte(testJSON)
	expectedData := []byte(testString)
	testPath := "testfile.bin"

	t.Run("successful file upload", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := tmpDir + "/" + testPath

		result, err := UploadFile(filePath, testContent)
		require.NoError(t, err)

		expectedMsg := fmt.Sprintf(i18n.BinInputSuccess+"%s\n", filePath)
		assert.Equal(t, expectedMsg, result)

		fileData, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, expectedData, fileData)
	})

	t.Run("invalid json content", func(t *testing.T) {
		invalidContent := []byte(`{"bin":123}`)
		_, err := UploadFile(testPath, invalidContent)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json unmarshal failed")
	})

	t.Run("invalid file path", func(t *testing.T) {
		invalidPath := "/invalid/path/to/file.bin"
		_, err := UploadFile(invalidPath, testContent)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file write failed")
	})

	t.Run("empty content", func(t *testing.T) {
		_, err := UploadFile(testPath, []byte{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json unmarshal failed")
	})

	t.Run("permission denied", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		tmpDir := t.TempDir()
		err := os.Chmod(tmpDir, 0555)
		require.NoError(t, err)
		defer os.Chmod(tmpDir, 0755)

		filePath := tmpDir + "/" + testPath
		_, err = UploadFile(filePath, testContent)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file write failed")
	})
}

func TestBinFileStruct(t *testing.T) {
	testString := "test content"
	encodedContent := base64.StdEncoding.EncodeToString([]byte(testString))
	content := []byte(`{"bin":"` + encodedContent + `"}`)

	var bf BinFile
	err := json.Unmarshal(content, &bf)
	require.NoError(t, err)
	assert.Equal(t, []byte(testString), bf.Data)
}
