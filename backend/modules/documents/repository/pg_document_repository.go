package repository

import (
	"backend/modules/documents/helpers/apperrors"
	"backend/modules/documents/models"

	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type pgDocumentRepository struct {
	DB *sqlx.DB
}

func NewDocumentRepository(db *sqlx.DB) models.DocumentRepository {
	return &pgDocumentRepository{
		DB: db,
	}
}

/******************** Core CRUD Functions ********************/

func (pgRepo *pgDocumentRepository) CreateDocumenteRecord(ctx context.Context, document *models.Document) (*uuid.UUID, error) {
	query := `INSERT INTO documents(document_id, document_title, author_name, author_id, description, cdn_url, security_access_level, created_at, updated_at, language) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING document_id`
	if err := pgRepo.DB.GetContext(ctx, document, query, document.DocumentID, document.DocumentTitle, document.AuthorName, document.AuthorID, document.Description,
		document.CDN_URL, document.ClassificationLevel, document.CreatedAt, document.UpdatedAt, document.Language); err != nil {
		log.Printf(`error: Could not create document: %v. Reason: %v`, document.DocumentTitle, err)
		return nil, apperrors.NewInternal()
	}
	return &document.DocumentID, nil
}

func (pgRepo *pgDocumentRepository) GetAllDocumentRecords(ctx context.Context) (*[]models.Document, error) {
	log.Print("Hello from pgRepo: ")
	queryResult := &[]models.Document{}
	//queryGetAllRecords := `SELECT * FROM documents;`
	//date_part('day', created_at) AS created_at,
	queryGetAllRecords := `select 
													document_id, 
													document_title, 
													author_name, 
													author_id,
													description, 
													security_access_level, 
													to_char(created_at, 'yyyy-mm-dd') AS created_at,
													to_char(created_at, 'yyyy-mm-dd') AS updated_at, 
													language 
												from documents;`

	if err := pgRepo.DB.SelectContext(ctx, queryResult, queryGetAllRecords); err != nil {
		log.Printf("error in retriving Document records from Postgres, error: %v, query: %v", err, queryGetAllRecords)
		return queryResult, apperrors.NewInternal()
	}
	return queryResult, nil
}

func (pgRepo *pgDocumentRepository) GetDocumentByID(ctx context.Context, documentID string) (*models.Document, error) {
	log.Print("Hello from pgRepo: ")
	queryResult := &models.Document{}
	queryGetRecord := `select 
													document_id, 
													document_title, 
													author_name, 
													author_id,
													description, 
													security_access_level, 
													to_char(created_at, 'yyyy-mm-dd') AS created_at,
													to_char(updated_at, 'yyyy-mm-dd') AS updated_at, 
													language 
												from documents
												WHERE document_id = $1;`

	if err := pgRepo.DB.GetContext(ctx, queryResult, queryGetRecord, documentID); err != nil {
		log.Printf("error in retriving Document records from Postgres, error: %v, query: %v\n", err, queryGetRecord)
		return queryResult, apperrors.NewInternal()
	}
	return queryResult, nil
}

func (pgRepo *pgDocumentRepository) GetDocumentByID_RowScan(ctx context.Context, documentID string) (*models.Document, error) {
	log.Print("Hello from pgRepo: ")
	queryResult := &models.Document{}
	//queryGetAllRecords := `SELECT * FROM documents;`
	//date_part('day', created_at) AS created_at,
	query := `select 
													document_id, 
													document_title, 
													author_name, 
													author_id,
													description, 
													security_access_level, 
													to_char(created_at, 'yyyy-mm-dd') AS created_at,
													to_char(created_at, 'yyyy-mm-dd') AS updated_at, 
													language 
												from documents
												WHERE document_id = $1;`

	row := pgRepo.DB.QueryRowContext(ctx, query, documentID)
	err := row.Scan(
		&queryResult.DocumentID,
		&queryResult.DocumentTitle,
		&queryResult.AuthorName,
		&queryResult.AuthorID,
		&queryResult.Description,
		&queryResult.ClassificationLevel,
		&queryResult.CreatedAt,
		&queryResult.UpdatedAt,
		&queryResult.Language,
	)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func (pgRepo *pgDocumentRepository) GetDocsByPagination(ctx context.Context, pageIndex int, pageCount int) (*[]models.Document, error) {
	result := &[]models.Document{}

	query := `SELECT * FROM documents ORDER BY created_at LIMIT $1 OFFSET $2;`

	err := pgRepo.DB.SelectContext(ctx, result, query, pageIndex, pageCount)
	if err != nil {
		log.Print("")
	}

	return result, nil
}

func (pgRepo *pgDocumentRepository) DeleteDocument(ctx context.Context, documentID uuid.UUID, accountID uuid.UUID) error {
	query := `DELETE FROM documents WHERE document_id = $1 AND author_id = $2 RETURNING document_id`

	err := pgRepo.DB.QueryRowContext(ctx, query, documentID, accountID)
	if err != nil {
		log.Printf(`error: Could not delete document: %v. Reason: %v`, documentID, err)
		return apperrors.NewAuthorization("Unable to delete document")
	}
	log.Print(err)
	return nil
}

func (pgRepo *pgDocumentRepository) UpdateDocument(ctx context.Context, updateDocRequest *models.Document, accountID uuid.UUID) (*models.Document, error) {
	//document := &models.Document{}
	log.Print(updateDocRequest)
	query := `UPDATE documents SET
							document_title = COALESCE(NULLIF($1, document_title)),
							author_name = COALESCE($2, author_name),
							description = COALESCE($3, description),
							security_access_level = COALESCE($4, security_access_level),
							language = COALESCE($5, language),
							updated_at = COALESCE($6, updated_at)
						WHERE document_id = $7 AND author_id = $8;
	 					`

	updatedDocument, err := pgRepo.DB.MustBegin().ExecContext(ctx, query, updateDocRequest.DocumentTitle, updateDocRequest.AuthorName, updateDocRequest.Description, updateDocRequest.ClassificationLevel,
		updateDocRequest.Language, updateDocRequest.UpdatedAt, updateDocRequest.DocumentID, accountID)
	if err != nil {
		log.Printf("error in updating Document records in Postgres: %v,\n result: %v", updatedDocument, err)
		return nil, apperrors.NewInternal()
	}

	log.Print(updatedDocument)

	return updateDocRequest, nil
}

/******************** Search Functions ********************/

func (pgRepo *pgDocumentRepository) SearchDocsByDateRange(ctx context.Context, documentRequest *models.DocumentSearchRequest) (*[]models.Document, error) {
	log.Print("Hello from pgRepo: ", documentRequest)
	queryResult := &[]models.Document{}
	queryGetAllRecords := `SELECT * FROM documents WHERE created_at BETWEEN SYMMETRIC $1 AND $2;`
	if err := pgRepo.DB.SelectContext(ctx, queryResult, queryGetAllRecords, documentRequest.DatesRangeStart, documentRequest.DatesRangeEnd); err != nil {
		log.Printf("error in retriving Document records from Postgres, error: %v, query: %v", err, queryGetAllRecords)
		return queryResult, apperrors.NewInternal()
	}
	return queryResult, nil
}

func (pgRepo *pgDocumentRepository) SearchDocsByTitle(ctx context.Context, documentRequest *models.DocumentSearchRequest) (*[]models.Document, error) {
	log.Print("Hello from pgRepo DocumentsSearch: ", documentRequest)
	queryResult := &[]models.Document{}
	querySearch := `
									SELECT * FROM documents
									WHERE document_title ILIKE '%' || $1 || '%';
									`
	// `
	// SELECT * FROM documents_record
	// WHERE document_title LIKE $1 AND author_name LIKE $2 AND security_access_level LIKE $3;
	// `
	//($3 IS NOT NULL AND $4 IS NOT NULL OR created_at_date BETWEEN SYMMETRIC cast($3 as date) AND cast($4 as date))
	// queryGetAllRecords := `
	// SELECT * FROM documents_record
	// WHERE (document_title = $1) OR author_name = $2 OR created_at_date BETWEEN SYMMETRIC cast($3 as date) AND cast($4 as date) OR security_access_level = $5);
	// `

	// `
	// SELECT * FROM documents_record
	// WHERE ($1 IS NULL OR document_title = $1)
	// AND ($2 IS NULL OR author_name = $2)
	// AND ($3 IS NULL OR $4 IS NULL OR OR created_at_date BETWEEN SYMMETRIC cast($3 as date) AND cast($4 as date)
	// AND ($5 IS NULL OR security_access_level = $5);
	// `

	if err := pgRepo.DB.SelectContext(ctx, queryResult, querySearch, documentRequest.SearchTags); err != nil {
		log.Printf("error in retriving Document records from Postgres: %v, query:\n %v", err, querySearch)
		return queryResult, apperrors.NewInternal()
	}
	log.Print(queryResult)
	return queryResult, nil
}

func (pgRepo *pgDocumentRepository) Debug_CreateFileRecordWithContent(ctx context.Context, document *models.DocumentWebForm) error {
	query := `INSERT INTO files(document_id, document_title, author_name, author_id, description, document_content, security_access_level, created_at, updated_at, language) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`
	if err := pgRepo.DB.GetContext(ctx, document, query, document.DocumentTitle, document.DocumentID, document.AuthorName, document.AuthorID, document.Description,
		document.DocumentContent, document.ClassificationLevel, document.CreatedAt, document.UpdatedAt, document.Language); err != nil {
		log.Printf(`error: Could not create document: %v. Reason: %v`, document.DocumentTitle, err)
		return apperrors.NewInternal()
	}

	log.Printf("Debug: files error checks passed. Review insert for correctness.\n")
	log.Printf("Data passed: document_id:%v, document_title:%v, author_name:%v, author_id:%v, description:%v, document_content:%v, security_access_level:%v, created_at:%v, updated_at:%v, language:%v",
		document.DocumentID, document.DocumentTitle, document.AuthorName, document.AuthorID, document.Description, document.DocumentContent, document.ClassificationLevel, document.CreatedAt, document.UpdatedAt, document.Language)

	return nil
}
