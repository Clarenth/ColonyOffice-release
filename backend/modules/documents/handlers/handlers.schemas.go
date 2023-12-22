package handlers

import (
	"mime/multipart"

	"github.com/google/uuid"
)

/*
Schemas are used to structure incoming requests for the server handlers to understand.
Other names include Data Transfer Objects(DTOs) in Java and C#, and payloads or requests in JavaScript.
*/
type updateRequest struct {
}

// docFormRequest is used to save document data sent by the client using the in-app Document writer
type docFormRequest struct {
	Title               string `json:"title"`
	Author              string `json:"author"`
	Description         string `json:"description"`
	Content             string `json:"content"`
	SecurityAccessLevel string `json:"security_access_level"`
	//SecurityAccessLevel *models.SecurityAccessLevel `json:"security_access_level"`
}

type docRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	//Author              string `json:"author"`
	Description         string `json:"description"`
	SecurityAccessLevel string `json:"security_access_level"`
	Language            string `json:"language"`
	//Files               map[string][]*multipart.FileHeader
	//SecurityAccessLevel *models.SecurityAccessLevel `json:"security_access_level"`
}

// Document models the Document structure for user uploaded files
type docUpdateRequest struct {
	DocumentTitle       string    `json:"title"`
	DocumentID          uuid.UUID `json:"document_id"`
	AuthorName          string    `json:"author"`
	AuthorID            uuid.UUID `json:"author_id"`
	Description         string    `json:"description"`
	ClassificationLevel string    `json:"security_access_level"`
	Language            string    `db:"language" json:"language"`
}

// docUploadRequest is used by docUpload handler for incoming client requests
// when documents are uploaded instead of wirtten using the app Rich Text Editor
type docUploadFileRequest struct {
	Title               string `json:"title"`
	Author              string `json:"author"`
	Description         string `json:"description"`
	SecurityAccessLevel string `json:"security_access_level"`
	Language            string `json:"language"`
	Files               map[string][]*multipart.FileHeader
	//SecurityAccessLevel *models.SecurityAccessLevel `json:"security_access_level"`
}

type docGetRequest_Single struct {
	SearchTags          string `json:"search_tags"`
	Author              string `json:"author"`
	Colony              string `json:"colony,omitempty"`
	DatesRangeStart     string `json:"dates_range_start"`
	DatesRangeEnd       string `json:"dates_range_end"`
	SecurityAccessLevel string `json:"security_access_level"`
	//SecurityAccessLevel *models.SecurityAccessLevel `json:"security_access_level"`
}
