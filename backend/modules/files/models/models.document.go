package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// Document models the Document structure for user uploaded files
type Document struct {
	DocumentTitle       string                  `db:"document_title" json:"document_title"`
	DocumentID          uuid.UUID               `db:"document_id" json:"document_id"`
	AuthorName          string                  `db:"author_name" json:"author_name"`
	AuthorID            uuid.UUID               `db:"author_id" json:"-"`
	Description         string                  `db:"description" json:"description"`
	File                string                  `db:"-" json:"-"`
	File2               *multipart.FileHeader   `db:"-" json:"-"`
	Files               []*multipart.FileHeader `db:"-" json:"-"`
	CDN_URL             string                  `db:"cdn_url" json:"cdn_url"`
	ClassificationLevel string                  `db:"security_access_level" json:"security_access_level"`
	//ClassificationLevel *SecurityAccessLevel `db:"security_access_level" json:"security_access_level"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Language  string    `db:"language" json:"language"`
}

type DocumentSearchRequest struct {
	SearchTags          string `json:"search_tags"`
	Author              string `json:"author"`
	Colony              string `json:"colony,omitempty"`
	DatesRangeStart     string `json:"dates_range_start"`
	DatesRangeEnd       string `json:"dates_range_end"`
	SecurityAccessLevel string `json:"security_access_level"`
	//SecurityAccessLevel *models.SecurityAccessLevel `json:"security_access_level"`
}

type DocumentWebForm struct {
	DocumentTitle       string    `db:"document_title" json:"file_name"`
	DocumentID          uuid.UUID `db:"document_id" json:"document_id"`
	AuthorName          string    `db:"author_name" json:"author_name"`
	AuthorID            uuid.UUID `db:"author_id" json:"-"`
	Description         string    `db:"description" json:"description"`
	DocumentContent     string    `db:"document_content" json:"document_content"`
	CDN_URL             string    `db:"cdn_url" json:"cdn_url"`
	ClassificationLevel string    `db:"security_access_level" json:"security_access_level"`
	//ClassificationLevel *SecurityAccessLevel `db:"security_access_level" json:"security_access_level"`
	CreatedAt time.Time `db:"created_at_date" json:"created_at_date"`
	UpdatedAt time.Time `db:"updated_at_date" json:"updated_at_date"`
	Language  string    `db:"language" json:"language"`
}
