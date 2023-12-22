package models

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentHandlers interface {
	FormFileData(ctx *gin.Context)
	CreateDocument(ctx *gin.Context)
	GetAllDocuments(ctx *gin.Context)
	GetDocumentByID(ctx *gin.Context)
	GetDocsByPagination(ctx *gin.Context)
	DeleteDoc(ctx *gin.Context)
	UpdateDoc(ctx *gin.Context)
	SearchDocs(ctx *gin.Context)
}

type DocumentService interface {
	WriteDocumentForm(ctx context.Context, docRequest *Document) (*uuid.UUID, error)
	//WriteFormToFile
	DeleteDocument(ctx context.Context, documentID uuid.UUID, accountID uuid.UUID) error
	GetAllDocumentRecords(ctx context.Context) (*[]Document, error)
	GetDocumentByID(ctx context.Context, documentID string) (*Document, error)
	GetDocsByPagination(ctx context.Context, pageIndex int, pageCount int) (*[]Document, error)
	SearchDocsByDateRange(ctx context.Context, docRequest *DocumentSearchRequest) (*[]Document, error)
	SearchDocsByTitle(ctx context.Context, request *DocumentSearchRequest) (*[]Document, error)
	UpdateDocument(ctx context.Context, request *Document, accountID uuid.UUID) (*Document, error)
}

type DocumentRepository interface {
	CreateDocumenteRecord(ctx context.Context, docRequest *Document) (*uuid.UUID, error)
	GetAllDocumentRecords(ctx context.Context) (*[]Document, error)
	GetDocumentByID(ctx context.Context, documentID string) (*Document, error)
	GetDocumentByID_RowScan(ctx context.Context, documentID string) (*Document, error)
	GetDocsByPagination(ctx context.Context, pageIndex int, pageCount int) (*[]Document, error)
	DeleteDocument(ctx context.Context, documentID uuid.UUID, accountID uuid.UUID) error
	SearchDocsByDateRange(ctx context.Context, document *DocumentSearchRequest) (*[]Document, error)
	SearchDocsByTitle(ctx context.Context, documentRequest *DocumentSearchRequest) (*[]Document, error)
	UpdateDocument(ctx context.Context, updateDocRequest *Document, accountID uuid.UUID) (*Document, error)
}
