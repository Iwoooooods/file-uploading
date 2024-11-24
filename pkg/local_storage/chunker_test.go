package localstorage

import (
	"testing"
)

func TestUploader(t *testing.T) {
	t.Run("test chunking", func(t *testing.T) {
		chunker := DefaultFileChunker{
			chunkSize: 10,
		}

		chunks, err := chunker.ChunkFile("tmp/short_text.txt")
		if err != nil {
			t.Fatalf("error chunking file: %v", err)
		}
		t.Logf("chunks: %v", chunks)
	})

	t.Run("test large file chunking", func(t *testing.T) {
		chunker := LargeFileChunker{
			chunkSize: 10,
		}
		chunks, err := chunker.ChunkFile("tmp/long_text.txt")
		if err != nil {
			t.Fatalf("error chunking file: %v", err)
		}
		t.Logf("chunks: %v", chunks)
	})
}
