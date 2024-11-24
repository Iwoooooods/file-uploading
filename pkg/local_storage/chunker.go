package localstorage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
)

// ChunkFile chunks a file into smaller chunks
// returns a list of chunk metadata of each chunk
// it chunks the file sequentially based on the chunk size
func (c *DefaultFileChunker) ChunkFile(filePath string) ([]ChunkMeta, error) {
	var chunks []ChunkMeta

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, c.chunkSize)
	index := 0

	for {
		// read the file into buffer
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file")
				break
			}
			return nil, err
		}

		// add chunk metadata to chunks
		chunk, err := createChunk(filePath, index, buffer, bytesRead)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)

		index++
	}

	return chunks, nil
}

// ChunkLargeFile chunks a large file into chunks simultaneously
// instead of reading file in a loop
// returns a list of chunk metadata of each chunk
func (c *LargeFileChunker) ChunkFile(filePath string) ([]ChunkMeta, error) {
	numOfWorkers := 4 // number of workers to process chunks
	var chunks []ChunkMeta
	var wg sync.WaitGroup
	var mu sync.Mutex

	// open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// get the file size to determine the number of chunks
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	numOfChunks := int(fileInfo.Size() / int64(c.chunkSize))
	// one more chunk if the file size is not a multiple of chunk size
	if fileInfo.Size()%int64(c.chunkSize) != 0 {
		numOfChunks++
	}

	// chunk channel to store processed chunk metadata
	chunkChan := make(chan ChunkMeta, numOfChunks)
	// error channel to store errors
	errorChan := make(chan error, numOfChunks)
	// index channel to store chunk indices
	indexChan := make(chan int, numOfChunks)

	// populate the index channel with chunk indices
	// distributing work among multiple workers
	for i := 0; i < numOfChunks; i++ {
		indexChan <- i
	}
	close(indexChan)

	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range indexChan {
				// calculate the offset and length of the chunk
				offset := int64(index * c.chunkSize)
				length := c.chunkSize

				// read the chunk from the file
				file.Seek(offset, io.SeekStart)
				buffer := make([]byte, length)
				bytesRead, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					errorChan <- err
					return
				}

				// if bytesRead is 0, it means the end of the file
				if bytesRead > 0 {
					chunk, err := createChunk(filePath, index, buffer, bytesRead)
					if err != nil {
						errorChan <- err
						return
					}

					mu.Lock()
					chunks = append(chunks, chunk)
					mu.Unlock()

					chunkChan <- chunk
				}
			}
		}()
	}

	// block until all workers are done
	go func() {
		wg.Wait()
		close(chunkChan)
		close(errorChan)
	}()

	// get errors from error channel
	for err := range errorChan {
		return nil, err
	}

	return chunks, nil

}

func createChunk(filePath string, index int, buffer []byte, bytesRead int) (ChunkMeta, error) {
	// create unique MD5 hash for the chunk
	hash := md5.Sum(buffer[:bytesRead])
	hashStr := hex.EncodeToString(hash[:])

	// chunkName to be like file.txt.chunk.0
	chunkName := fmt.Sprintf("%s.chunk.%d", filePath, index)
	// create chunk file
	chunkFile, err := os.Create(chunkName)
	if err != nil {
		return ChunkMeta{}, err
	}
	defer chunkFile.Close()

	// write the chunk to the file
	_, err = chunkFile.Write(buffer[:bytesRead])
	if err != nil {
		return ChunkMeta{}, err
	}
	return ChunkMeta{
		FileName: chunkName,
		MD5Hash:  hashStr,
		Index:    index,
	}, nil
}
