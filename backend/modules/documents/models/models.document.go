package models

import (
	"time"

	"github.com/google/uuid"
)

// Document models the Document structure for user uploaded files
type Document struct {
	DocumentTitle       string    `db:"document_title" json:"document_title"`
	DocumentID          uuid.UUID `db:"document_id" json:"document_id"`
	AuthorName          string    `db:"author_name" json:"author_name"`
	AuthorID            uuid.UUID `db:"author_id" json:"author_id"`
	Description         string    `db:"description" json:"description"`
	CDN_URL             string    `db:"cdn_url" json:"cdn_url,omitempty"`
	ClassificationLevel string    `db:"security_access_level" json:"security_access_level"`
	//ClassificationLevel *SecurityAccessLevel `db:"security_access_level" json:"security_access_level"`

	CreatedAt string `db:"created_at" json:"created_at"` //time.Time look into converting all getDocument SQL calls to a time.Format. How expensive might this be
	UpdatedAt string `db:"updated_at" json:"updated_at"` //time.Time look into converting all getDocument SQL calls to a time.Format. How expensive might this be
	Language  string `db:"language" json:"language"`
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
