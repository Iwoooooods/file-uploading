package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Iwoooooods/fs-upload-go/internal/models"
)

type MetaRepositorySQLite struct {
	db *sql.DB
}

func NewMetaRepositorySQLite(db *sql.DB) *MetaRepositorySQLite {
	return &MetaRepositorySQLite{db}
}

func (r *MetaRepositorySQLite) Create(ctx context.Context, metadata models.FileMetadata) error {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO metadata (file_id, file_name, md5_hash) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, metadata.FileId, metadata.FileName, metadata.MD5Hash)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	log.Println("File metadata created successfully")

	return nil
}

func (r *MetaRepositorySQLite) Get(ctx context.Context, field string, value string) (models.FileMetadata, error) {
	stmt, err := r.db.PrepareContext(ctx, fmt.Sprintf("SELECT * FROM metadata WHERE %s = ?", field))
	if err != nil {
		log.Printf("failed to prepare statement: %v", err)
		return models.FileMetadata{}, err
	}
	defer stmt.Close()

	var metadata models.FileMetadata
	var id int64

	err = stmt.QueryRowContext(ctx, value).Scan(&id, &metadata.FileId, &metadata.FileName, &metadata.MD5Hash)
	if err != nil {
		return models.FileMetadata{}, err
	}

	return metadata, nil
}

func (r *MetaRepositorySQLite) Update(ctx context.Context, metadata models.FileMetadata) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE metadata SET file_name = ?, md5_hash = ? WHERE file_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, metadata.FileName, metadata.MD5Hash, metadata.FileId)
	if err != nil {
		return err
	}
	log.Println("File metadata updated successfully")

	return nil
}

func (r *MetaRepositorySQLite) Delete(ctx context.Context, fileId string) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM metadata WHERE file_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, fileId)
	if err != nil {
		return err
	}
	log.Println("File metadata deleted successfully")

	return nil
}
