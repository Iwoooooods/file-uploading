package localstorage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetadataManager(t *testing.T) {
	uploader, err := NewUploader(".", "tester")

	t.Run("new uploader entries", func(t *testing.T) {
		if err != nil {
			t.Errorf("failed to create uploader: %v", err)
		}
		if uploader.basePath != "./tester/" {
			t.Errorf("expected basePath to be ./tester/, got %v", uploader.basePath)
		}
		// user's folder should exist
		if _, err := os.Stat(uploader.basePath); os.IsNotExist(err) {
			t.Errorf("expected basePath to exist, got %v", err)
		}
		// user's metadata.json should exist
		if _, err := os.Stat(filepath.Join(uploader.basePath, "metadata.json")); os.IsNotExist(err) {
			t.Errorf("expected metadata.json to exist, got %v", err)
		}
	})
	t.Run("save metadata", func(t *testing.T) {
		fileName := "testfile"
		md5Hash := "md5_hash"
		metadata := make(map[string]FileMetadata)
		metadata[fileName] = FileMetadata{
			FileName: fileName,
			MD5Hash:  md5Hash,
		}
		t.Logf("metadata: %v", uploader.metaManager)
		err := uploader.metaManager.SaveMetadata(uploader.basePath, metadata)
		if err != nil {
			t.Errorf("failed to save metadata: %v", err)
		}
	})
	t.Run("load metadata", func(t *testing.T) {
		metadata, err := uploader.metaManager.LoadMetadata(uploader.basePath)
		if err != nil {
			t.Errorf("failed to load metadata: %v", err)
		}
		if len(metadata) != 1 {
			t.Errorf("expected only 1 metadata for one file, got %d", len(metadata))
		}
		for _, fileMetadata := range metadata {
			t.Logf("file_name: %v", fileMetadata.FileName)
			t.Logf("md5_hash: %v", fileMetadata.MD5Hash)
		}
	})

}
