package localstorage

import (
	"log"
	"os"
	"path/filepath"
)

func NewUploader(serverURL string, username string) (*DefaultUploader, error) {
	// check if user exists, if not create new fold for the user
	basePath := serverURL + "/" + username + "/"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// give permission to the user to read and write to the folder
		os.MkdirAll(basePath, os.ModePerm)
	}

	metaManager, err := NewMetadataManager(basePath)
	if err != nil {
		log.Printf("failed to create metadata manager: %v", err)
		return nil, err
	}

	return &DefaultUploader{
		serverURL:   serverURL,
		username:    username,
		basePath:    basePath,
		metaManager: metaManager,
	}, nil
}

// func (u *DefaultUploader) UploadFile(filePath string, md5 string) (fileId string, err error) {
// 	return nil
// }

func (u *DefaultUploader) CheckFileExists(fileName string, md5 string) (exists bool, err error) {
	// check if file exists
	if _, err := os.Stat(filepath.Join(u.basePath, fileName)); os.IsNotExist(err) {
		// if not exists, create new metadata.json and return fileId and false
		newMetadata := make(map[string]FileMetadata)
		newMetadata[fileName] = FileMetadata{
			FileName: fileName,
			MD5Hash:  md5,
		}

		err = u.metaManager.SaveMetadata(u.basePath, newMetadata)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// if exists, check if md5 changes
	// if md5 changes, update metadata.json
	// if md5 does not change, directly return fileId
	metadataFilePath := filepath.Join(u.basePath, "metadata.json")
	oldMetadata, err := u.metaManager.LoadMetadata(metadataFilePath)
	if err != nil {
		return false, err
	}

	for _, fileMetadata := range oldMetadata {
		if fileMetadata.MD5Hash == md5 {
			// file exists
			return true, nil
		}
	}

	return false, nil
}
