package localstorage

// metadata is stored in a json file
// when re-chunking, used to compare
type ChunkMeta struct {
	// when encoding to json, map the key to this struct field name
	FileName string `json:"file_name"`
	MD5Hash  string `json:"md5_hash"`
	Index    int    `json:"index"`
}

// metadata is stored in a json file
// example(test/testusername/testfile/metadata.json):
//
//	{
//		"example_name.md5hash.txt": {
//			"file_name": "example_name.md5hash.txt",
//			"md5_hash": "md5_hash"
//		}
//	}
//
// plus, the file id consists of filename+md5hash+extension
type FileMetadata struct {
	FileName string `json:"file_name"`
	MD5Hash  string `json:"md5_hash"`
}

type Config struct {
	ChunkSize int
	ServerURL string
}

type Chunker interface {
	ChunkFile(filePath string) ([]ChunkMeta, error)
	// ChunkLargeFile(filePath string) ([]ChunkMeta, error)
}

type DefaultFileChunker struct {
	chunkSize int
}

type LargeFileChunker struct {
	chunkSize int
}

type Uploader interface {
	UploadFile(filePath string, md5 string, metaManager MetadataManager) (fileId string, err error)
}

type DefaultUploader struct {
	serverURL   string
	username    string
	basePath    string
	metaManager *DefaultMetadataManager
}

// metadata is stored in a json file
type MetadataManager interface {
	SaveMetadata(filePath string, metadata map[string]FileMetadata) error
	LoadMetadata(filePath string) (map[string]FileMetadata, error)
}

type DefaultMetadataManager struct {
	metaData map[string]FileMetadata
}
