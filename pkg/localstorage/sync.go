package localstorage

import (
	"sync"
)

func synchronize(chunks []ChunkMeta, wg *sync.WaitGroup, metadata map[string]ChunkMeta, mu *sync.Mutex) {
	numOfWorkers := 4
	var chunkChan = make(chan ChunkMeta, len(chunks))
	var errorChan = make(chan error, len(chunks))

	for _, chunk := range chunks {
		chunkChan <- chunk
	}
	close(chunkChan)

	for i := 0; i < numOfWorkers; i++ {
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				newHash := chunk.MD5Hash

				// check if the chunk exists in the database
				oldChunk, ok := metadata[chunk.FileName]
				if !ok || oldChunk.MD5Hash != newHash {
					// TODO: update the chunk metadata
					metadata[chunk.FileName] = chunk
				}

				mu.Lock()
				metadata[chunk.FileName] = chunk
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	close(errorChan)
}
