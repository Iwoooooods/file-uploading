package blobstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

const (
	CHUNK_SIZE = 1024 * 1024 * 4 // 4MB
)

type File struct {
	Id   string
	Path string
	Ext  string
}

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
	Upload(ctx context.Context, data io.Reader, fileName string) (string, error)
	Read(ctx context.Context, fileId string) (io.Reader, error)
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

func (fs *LocalFS) Upload(ctx context.Context, data io.Reader, fileName string) (File, error) {
	fileId := strings.Split(fileName, ".")[0] + "-" + uuid.New().String()
	filePath := fileId + "." + strings.Split(fileName, ".")[1]
	chunks := make([]*Chunk, 0)

	buffer := make([]byte, CHUNK_SIZE)
	chunkId := int64(0)

	chunkPath := filepath.Join(fs.basePath, fileId)
	if err := os.MkdirAll(chunkPath, 0755); err != nil {
		return File{}, err
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
			return File{}, err
		}

		chunkPath := filepath.Join(fs.basePath, fileId, fmt.Sprintf("%d.chunk", chunkId))

		if err := fs.storeChunk(chunkPath, buffer[:n]); err != nil {
			return File{}, err
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
	return File{
		Id:   fileId,
		Path: filePath,
		Ext:  strings.Split(fileName, ".")[1],
	}, nil
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

func (fs *LocalFS) Read(ctx context.Context, downloadPath string) (io.Reader, error) {
	chunkParentPath := filepath.Join(fs.basePath, strings.Split(downloadPath, ".")[0])

	fmt.Println("reading chunks from: ", chunkParentPath)
	if _, err := os.Stat(chunkParentPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found")
	}

	chunkFiles, err := os.ReadDir(chunkParentPath)
	if err != nil {
		return nil, err
	}
	fmt.Println("found ", len(chunkFiles), " chunks")

	var buffer bytes.Buffer
	for _, chunk := range chunkFiles {
		fmt.Println("reading chunk: ", chunk.Name())
		chunkPath := filepath.Join(chunkParentPath, chunk.Name())
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return nil, err
		}
		buffer.Write(chunkData)
	}

	return &buffer, nil
}
