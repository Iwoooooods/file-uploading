package repositories

import (
	"context"

	"github.com/Iwoooooods/fs-upload-go/internal/models"
)

type MetaRepository interface {
	Create(ctx context.Context, metadata models.FileMetadata) error
	Get(ctx context.Context, field string, value string) (models.FileMetadata, error)
	Update(ctx context.Context, metadata models.FileMetadata) error
	Delete(ctx context.Context, fileId string) error
}
