package repository

import (
	"backend/modules/files/helpers/apperrors"
	"backend/modules/files/models"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type pgFilesRepository struct {
	DB *sqlx.DB
}

func NewFilesRepository(db *sqlx.DB) models.FilesRepository {
	return &pgFilesRepository{
		DB: db,
	}
}

func (pgRepo *pgFilesRepository) CreateFileRecord(ctx context.Context, file *models.File) error {
	log.Print("Hello from files PGRepo")
	query := `INSERT INTO files(file_id, document_id, file_title, title_hash, author_id, author_name, security_access_level, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *`

	if err := pgRepo.DB.GetContext(ctx, file, query, file.FileID, file.DocumentID, file.Title, file.TitleHash, file.AuthorID, file.AuthorName, file.SecurityAccessLevel,
		file.CreatedAt, file.UpdatedAt); err != nil {
		log.Printf(`error: Could not create file: %v. Reason: %v`, file.Title, err)
		return apperrors.NewInternal()
	}

	log.Printf("Debug: files_db_record error checks passed. Review insert for correctness.\n")
	log.Printf("Data passed: document_title:%v, document_id:%v, author_id:%v, author_name:%v, security_access_level:%v, created_at:%v, updated_at:%v",
		file.Title, file.FileID, file.AuthorID, file.AuthorName, file.SecurityAccessLevel, file.CreatedAt, file.UpdatedAt)
	return nil
}

func (pgRepo *pgFilesRepository) DeleteFile(ctx context.Context, filesID uuid.UUID, accountID uuid.UUID) error {
	query := `DELETE FROM files WHERE author_id = $1 AND files_id $2`

	err := pgRepo.DB.QueryRowContext(ctx, query, filesID, accountID)
	if err != nil {
		log.Print("error with deleting file from Postgres")
	}

	return nil
}

func (pgRepo *pgFilesRepository) GetFileByDocID(ctx context.Context, documentID uuid.UUID) (*[]models.File, error) {
	files := &[]models.File{}

	// query := `SELECT * from files WHERE document_id = $1`
	query := `select
							file_id,
							document_id,
							file_title,
							author_name,
							author_id,
							security_access_level,
							to_char(created_at, 'yyyy-mm-dd') AS created_at,
							to_char(updated_at, 'yyyy-mm-dd') AS updated_at
						from files WHERE document_id = $1;`

	if err := pgRepo.DB.SelectContext(ctx, files, query, documentID); err != nil {
		log.Printf("error in retriving Postgres rows. Error: %v, query: %v", err, query)
		return nil, apperrors.NewBadRequest("Could not find document")
	}
	log.Print(files)
	return files, nil
}
