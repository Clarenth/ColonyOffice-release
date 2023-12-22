package models

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FilesHandlers interface {
	UploadFiles(ctx *gin.Context)
	GetFile(ctx *gin.Context)
	GetFilesByDocID(ctx *gin.Context)
	DeleteFile(ctx *gin.Context)
	UpdateFile(ctx *gin.Context)
}

type FilesService interface {
	UploadFile(ctx context.Context, fileData *File) error
	DeleteFile(ctx context.Context, filesID uuid.UUID, accountID uuid.UUID) error
	GetFileByDocID(ctx context.Context, documentID uuid.UUID) (*[]File, error)
	UpdateFiles(ctx context.Context, filesID uuid.UUID, accountID uuid.UUID) error
}

type FilesRepository interface {
	CreateFileRecord(ctx context.Context, filesData *File) error
	DeleteFile(ctx context.Context, documentID uuid.UUID, accountID uuid.UUID) error
	GetFileByDocID(ctx context.Context, documentID uuid.UUID) (*[]File, error)
}
