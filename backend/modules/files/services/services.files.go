package services

import (
	"backend/modules/files/models"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// doountService acts as a struct for injecting an implementation of
// EmployeeRepository for use in the service methods
type filesService struct {
	FilesRepository models.FilesRepository
}

// ConfigAccountService will hold repositories that will eventually be injected
// into this service layer
type ConfigFilesService struct {
	FilesRepository models.FilesRepository
}

// NewAccountService is a factory function for initializing an employeeService
// with it's repository layer dependencies
func NewFilesService(ctx *ConfigFilesService) models.FilesService {
	return &filesService{
		FilesRepository: ctx.FilesRepository,
	}
}

func (services *filesService) UploadFile(ctx context.Context, fileData *models.File) error {
	log.Print("Hello from services.UplaodFile")

	fileData.CreatedAt = time.Now().UTC().String()
	fileData.UpdatedAt = fileData.CreatedAt

	err := services.FilesRepository.CreateFileRecord(ctx, fileData)
	if err != nil {
		return err
	}
	return nil
}

func (services *filesService) DeleteFile(ctx context.Context, docuemntID uuid.UUID, accountID uuid.UUID) error {
	//panic("not done yet")
	err := services.FilesRepository.DeleteFile(ctx, docuemntID, accountID)
	if err != nil {
		log.Print("error with deleting from Postgres")
	}

	return nil
}

func (services *filesService) GetFileByDocID(ctx context.Context, documentID uuid.UUID) (*[]models.File, error) {
	result, err := services.FilesRepository.GetFileByDocID(ctx, documentID)
	if err != nil {
		log.Print("erorr in services GetFileByDocID. Error: ", err)
		return nil, err
	}
	return result, nil
}

func (service *filesService) UpdateFiles(ctx context.Context, filesID uuid.UUID, accountID uuid.UUID) error {
	panic("not done yet")
}
