package blobstorage

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalFS_Upload(t *testing.T) {
	fs := NewLocalFS("./tmp/")
	t.Run("Upload small file", func(t *testing.T) {
		fileId, err := fs.Upload(context.Background(), bytes.NewReader([]byte("hello, world")))
		assert.NoError(t, err)
		assert.NotEmpty(t, fileId)
	})

	// Test uploading a larger file that will be split into chunks
	t.Run("Upload large file", func(t *testing.T) {
		// Create test data larger than chunk size
		largeData := make([]byte, CHUNK_SIZE*2+1024) // 2 chunks + 1KB
		for i := range largeData {
			largeData[i] = byte(i % 256) // Fill with repeating pattern
		}

		fileId, err := fs.Upload(context.Background(), bytes.NewReader(largeData))
		assert.NoError(t, err)
		assert.NotEmpty(t, fileId)
	})
}
