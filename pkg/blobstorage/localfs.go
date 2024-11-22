package blobstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

const (
	CHUNK_SIZE = 1024 * 1024 * 4 // 4MB
)

type Chunk struct {
	data      io.Reader
	chunkId   int64
	chunkPath string
	next      *Chunk
}

func NewChunk(data io.Reader, chunkId int64, chunkPath string) *Chunk {
	return &Chunk{
		data:      data,
		chunkId:   chunkId,
		chunkPath: chunkPath,
	}
}

type BlobStorage interface {
}

type LocalFS struct {
	basePath string
	mu       sync.Mutex
}

func NewLocalFS(basePath string) *LocalFS {
	return &LocalFS{
		basePath: basePath,
	}
}

func (fs *LocalFS) Upload(ctx context.Context, data io.Reader) (string, error) {
	fileId := uuid.New().String()
	chunks := make([]*Chunk, 0)

	buffer := make([]byte, CHUNK_SIZE)
	chunkId := int64(0)

	chunkPath := filepath.Join(fs.basePath, fileId)
	if err := os.MkdirAll(chunkPath, 0755); err != nil {
		return "", err
	}

	for {
		n, err := io.ReadAtLeast(data, buffer, 1)
		if err == io.EOF {
			break
		}
		if n == 0 {
			break
		}
		if err != nil {
			return "", err
		}

		chunkPath := filepath.Join(fs.basePath, fileId, fmt.Sprintf("%d.chunk", chunkId))

		if err := fs.storeChunk(chunkPath, buffer[:n]); err != nil {
			return "", err
		}

		chunk := NewChunk(bytes.NewReader(buffer[:n]), chunkId, chunkPath)
		if len(chunks) > 0 {
			chunks[len(chunks)-1].next = chunk
		}
		chunks = append(chunks, chunk)
		chunkId++

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			fmt.Printf("Uploaded %d chunks\n", chunkId)
			break
		}
	}
	return fileId, nil
}

func (fs *LocalFS) storeChunk(chunkPath string, data []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	f, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)

	return err
}
