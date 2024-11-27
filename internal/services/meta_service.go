package services

import (
	"context"

	"github.com/Iwoooooods/fs-upload-go/internal/models"
	"github.com/Iwoooooods/fs-upload-go/internal/repositories"
)

type MetaService interface {
	SaveMetadata(ctx context.Context, metadata models.FileMetadata) error
	GetMetadataById(ctx context.Context, fileId string) (models.FileMetadata, error)
	GetMetadataByMD5(ctx context.Context, md5 string) (models.FileMetadata, error)
	UpdateMetadata(ctx context.Context, metadata models.FileMetadata) error
	DeleteMetadata(ctx context.Context, fileId string) error
}

type MetaServiceImpl struct {
	repo repositories.MetaRepository
}

func NewMetaService(repo repositories.MetaRepository) *MetaServiceImpl {
	return &MetaServiceImpl{repo}
}

// SaveMetadata saves file metadata to the database
func (s *MetaServiceImpl) SaveMetadata(ctx context.Context, metadata models.FileMetadata) error {
	return s.repo.Create(ctx, metadata)
}

// GetMetadata retrieves file metadata from the database
func (s *MetaServiceImpl) GetMetadataById(ctx context.Context, fileId string) (models.FileMetadata, error) {
	return s.repo.Get(ctx, "file_id", fileId)
}

func (s *MetaServiceImpl) GetMetadataByMD5(ctx context.Context, md5 string) (models.FileMetadata, error) {
	return s.repo.Get(ctx, "md5_hash", md5)
}

// UpdateMetadata updates existing file metadata in the database
func (s *MetaServiceImpl) UpdateMetadata(ctx context.Context, metadata models.FileMetadata) error {
	return s.repo.Update(ctx, metadata)
}

// DeleteMetadata removes file metadata from the database
func (s *MetaServiceImpl) DeleteMetadata(ctx context.Context, fileId string) error {
	return s.repo.Delete(ctx, fileId)
}
