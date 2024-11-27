package models

type FileMetadata struct {
	FileId   string `json:"file_id" db:"file_id"`
	FileName string `json:"file_name" db:"file_name"`
	MD5Hash  string `json:"md5_hash" db:"md5_hash"`
}
