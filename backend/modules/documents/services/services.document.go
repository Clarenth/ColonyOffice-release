package services

import (
	"backend/modules/documents/models"

	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// doountService acts as a struct for injecting an implementation of
// EmployeeRepository for use in the service methods
type documentService struct {
	DocumentDirectory  string
	DocumentRepository models.DocumentRepository
}

// ConfigAccountService will hold repositories that will eventually be injected
// into this service layer
type ConfigDocumentService struct {
	DocumentDirectory  string
	DocumentRepository models.DocumentRepository
}

// NewAccountService is a factory function for initializing an employeeService
// with it's repository layer dependencies
func NewDocumentService(ctx *ConfigDocumentService) models.DocumentService {
	return &documentService{
		DocumentDirectory:  ctx.DocumentDirectory,
		DocumentRepository: ctx.DocumentRepository,
	}
}

/******************** Core CRUD Functions ********************/
func (services *documentService) WriteDocumentForm(ctx context.Context, docData *models.Document) (*uuid.UUID, error) {
	log.Print("Hello from services.WriteDocToFile")

	docData.DocumentID = uuid.New()
	docData.CreatedAt = time.Now().UTC().String()
	docData.UpdatedAt = docData.CreatedAt
	log.Print(docData)

	document, err := services.DocumentRepository.CreateDocumenteRecord(ctx, docData)
	if err != nil {
		return nil, err
	}
	return document, err
}

func (services *documentService) GetAllDocumentRecords(ctx context.Context) (*[]models.Document, error) {
	log.Print("Hello from services: ")

	docRecords, err := services.DocumentRepository.GetAllDocumentRecords(ctx)
	if err != nil {
		log.Print("error retriving document records for request")
		return nil, err
	}
	return docRecords, nil
}

func (services *documentService) GetDocumentByID(ctx context.Context, documentID string) (*models.Document, error) {
	log.Print("Hello from services")

	document, err := services.DocumentRepository.GetDocumentByID_RowScan(ctx, documentID)
	if err != nil {
		log.Print("error with retriving SQL record from pgRepo", err)
	}

	return document, nil
}

func (services *documentService) GetDocsByPagination(ctx context.Context, pageIndex int, pageCount int) (*[]models.Document, error) {
	log.Print("hello from services")

	result, err := services.DocumentRepository.GetDocsByPagination(ctx, pageIndex, pageCount)
	if err != nil {
		log.Print("error in services GetDocsByPagination")
		return nil, err
	}

	return result, nil
}

func (services *documentService) DeleteDocument(ctx context.Context, documentID uuid.UUID, accountID uuid.UUID) error {
	err := services.DocumentRepository.DeleteDocument(ctx, documentID, accountID)
	if err != nil {
		log.Print("error with deleting the document ")
		return err
	}

	return nil
}

func (services *documentService) UpdateDocument(ctx context.Context, request *models.Document, accountID uuid.UUID) (*models.Document, error) {
	request.UpdatedAt = time.Now().UTC().String()

	updatedDoc, err := services.DocumentRepository.UpdateDocument(ctx, request, accountID)
	if err != nil {
		log.Print("error with Postgres during Document update")
		return nil, err
	}

	return updatedDoc, nil
}

/******************** Search Functions ********************/

func (services *documentService) SearchDocsByDateRange(ctx context.Context, request *models.DocumentSearchRequest) (*[]models.Document, error) {
	log.Print("Hello from services: ", request)

	docRecords, err := services.DocumentRepository.SearchDocsByTitle(ctx, request)
	if err != nil {
		log.Print("error retriving document records for request")
		return nil, err
	}
	return docRecords, nil
}

func (services *documentService) SearchDocsByTitle(ctx context.Context, request *models.DocumentSearchRequest) (*[]models.Document, error) {
	log.Print("Hello from services: ", request)

	docRecords, err := services.DocumentRepository.SearchDocsByTitle(ctx, request)
	if err != nil {
		log.Print("error retriving document records for request")
		return nil, err
	}
	return docRecords, nil
}

// func (services *documentService) _writeUploadToFile(ctx context.Context, docData *models.Document) error {
// 	log.Print("Hello from services.WriteUploadToFile")

// 	docData.DocumentID = uuid.New()
// 	docData.CreatedAt = time.Now().UTC()
// 	docData.UpdatedAt = docData.CreatedAt

// 	fileName := strings.ReplaceAll(docData.CreatedAt.String()+"."+docData.File2.Filename, " ", "_")
// 	docData.File2.Filename = fileName
// 	docData.CDN_URL = services.DocumentDirectory + "/" + fileName
// 	log.Print(fileName)

// 	fileOutput, err := docData.File2.Open()
// 	//fileOutput, err := os.OpenFile(filepath.Join(services.DocumentDirectory, filepath.Base(docData.DocumentTitle)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Printf("erorr: document directory does not exist, or some other error with file path.")
// 		return apperrors.NewInternal()
// 	}
// 	defer fileOutput.Close()

// 	dst, err := os.Create(docData.CDN_URL)
// 	if err != nil {
// 		log.Print("error in creating docData.CDN_URL")
// 		return err
// 	}
// 	defer dst.Close()

// 	_, err = io.Copy(dst, fileOutput)
// 	if err != nil {
// 		log.Print("error in io.Copy of file upload")
// 	}

// 	// if err := ctx.SaveUploadedFile(docData.File, docData.CDN_URL); err != nil {
// 	// 	log.Print(docData.File)
// 	// 	log.Print("error in fileOutput.Write()")
// 	// 	return err
// 	// }

// 	services.DocumentRepository.CreateUploadFileRecord(ctx, docData)
// 	return nil
// }
