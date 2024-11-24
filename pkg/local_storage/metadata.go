package localstorage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func NewMetadataManager(basePath string) (*DefaultMetadataManager, error) {
	// read metadata.json of the uploader
	metadata, err := os.ReadFile(filepath.Join(basePath, "metadata.json"))
	// if user does not have metadata.json, create an empty one
	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(basePath, "metadata.json"), []byte("{}"), os.ModePerm)
		return &DefaultMetadataManager{
			metaData: make(map[string]FileMetadata),
		}, nil
	}
	if err != nil && err != os.ErrNotExist {
		return nil, err
	}

	var metadataMap map[string]FileMetadata
	err = json.Unmarshal(metadata, &metadataMap)
	if err != nil {
		return nil, err
	}

	return &DefaultMetadataManager{
		metaData: metadataMap,
	}, nil
}

// SaveMetadata saves the metadata to the given file path
// metaPath is defined: basePath/username/metadata.json
// metadata is a map of file name to FileMetadata
func (m *DefaultMetadataManager) SaveMetadata(metaPath string, metadata map[string]FileMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(metaPath, "metadata.json"), data, 0644)
}

// LoadMetadata loads the metadata from the given file path
// metaPath is defined: basePath/username/metadata.json
func (m *DefaultMetadataManager) LoadMetadata(basePath string) (map[string]FileMetadata, error) {
	metaPath := filepath.Join(basePath, "metadata.json")
	metadata := make(map[string]FileMetadata)

	data, err := os.ReadFile(metaPath)
	if err != nil {
		return metadata, err
	}

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}
