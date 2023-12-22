package models

import (
	"github.com/google/uuid"
)

type File struct {
	FileID              uuid.UUID `db:"file_id" json:"file_id"`
	DocumentID          uuid.UUID `db:"document_id" json:"document_id"`
	Title               string    `db:"file_title" json:"file_title"`
	TitleHash           string    `db:"title_hash" json:"-"`
	AuthorID            uuid.UUID `db:"author_id" json:"author_id"`
	AuthorName          string    `db:"author_name" json:"author_name"`
	SecurityAccessLevel string    `db:"security_access_level" json:"security_access_level"`
	CreatedAt           string    `db:"created_at" json:"created_at"`
	UpdatedAt           string    `db:"updated_at" json:"updated_at"`
	// CreatedAt           time.Time `db:"created_at" json:"created_at"`
	// UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}
