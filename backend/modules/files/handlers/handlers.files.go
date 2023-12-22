package handlers

import (
	auth "backend/modules/auth/models"
	"backend/modules/files/helpers/apperrors"
	"backend/modules/files/helpers/hashing"
	"backend/modules/files/models"

	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler holds the methods expected by the service layer for handling routing
type handler struct {
	DocumentDirectory string
	FilesService      models.FilesService
}

// HandlerConfig holds the server configuration need by by the handler
type HandlerConfig struct {
	DocumentDirectory string
	FilesService      models.FilesService
}

func NewFilesHandlers(ctx *HandlerConfig) models.FilesHandlers {
	return &handler{
		DocumentDirectory: ctx.DocumentDirectory,
		FilesService:      ctx.FilesService,
	}
}

func (handler *handler) UploadFiles(ctx *gin.Context) {
	authHeader := ctx.MustGet("account").(*auth.JWTToken)

	fileMetadata := &models.File{
		AuthorID: authHeader.IDCode,
	}

	// define some variables used throughout the function
	// n: for keeping track of bytes read and written
	// err: for storing errors that need checking
	var n int
	var err error

	// define pointers for the multipart reader and its parts
	var mr *multipart.Reader
	var part *multipart.Part
	//formWatcher := part.Header.Get("Content-Type")

	log.Println("File Upload Endpoint Hit")

	if mr, err = ctx.Request.MultipartReader(); err != nil {
		log.Print("error with upload. We should not see this.")
		log.Printf("Hit error while opening multipart reader: %s", err.Error())
		ctx.JSON(apperrors.NewInternal().Status(), gin.H{
			"error": err,
		})
		return
	}

	// buffer to be used for reading bytes from files
	chunk := make([]byte, 4096)

	// continue looping through all parts, *multipart.Reader.NextPart() will
	// return an End of File when all parts have been read.
	for {
		// variables used in this loop only
		// tempfile: filehandler for the temporary file
		// filesize: how many bytes where written to the tempfile
		// uploaded: boolean to flip when the end of a part is reached
		//var tempfile *os.File
		var filesize int
		var uploaded bool
		var destfile *os.File

		if part, err = mr.NextPart(); err != nil {
			if err != io.EOF {
				log.Printf("Hit error while fetching next part: %s", err.Error())
				ctx.JSON(apperrors.NewInternal().Status(), gin.H{
					//fmt.Fprintf(w, "Error occured during upload"),
					"error": err,
				})
			} else {
				log.Printf("Hit last part of multipart upload")
				ctx.JSON(http.StatusOK, gin.H{
					"result": "file saved successfully",
				})
			}
			return
		}
		// at this point the filename and the mimetype is known
		log.Printf("Uploaded filename: %s", part.FileName())
		log.Printf("Uploaded mimetype: %s", part.Header)

		if part.FormName() == "document_id" {
			bytes, err := io.ReadAll(part)
			if err != nil {
				log.Printf("error with reading document id form: %v", err)
			}
			formValue := string(bytes)
			fileMetadata.DocumentID = uuid.MustParse(formValue)
			continue
		}

		if part.FormName() == "security_access_level" {
			j, _ := io.ReadAll(part)
			log.Print(j)
			k := string(j)
			log.Print(k)
			fileMetadata.SecurityAccessLevel = k
			continue
		}
		fileMetadata.FileID = uuid.New()
		fileMetadata.Title = part.FileName()

		hash := hashing.HashFile(fileMetadata.Title)
		location := filepath.Join(handler.DocumentDirectory, hash.PathName)
		if err := os.MkdirAll(location, os.ModePerm); err != nil {
			log.Print(err)
			return
		}

		// if part.Header.Get("Content-Type") == "application/json" {
		// 	log.Print("Hello application/json")
		// 	uploadMetadata, err := part.Read(chunk[:])
		// 	if err != nil {
		// 		log.Print(uploadMetadata)
		// 		log.Print(err)
		// 		return
		// 	}
		// 	log.Print(uploadMetadata)
		// }

		destfile, err = os.Create(filepath.Join(handler.DocumentDirectory, hash.Filename()))

		//tempfile, err = os.CreateTemp(os.TempDir(), "example-upload-*.tmp")
		if err != nil {
			log.Printf("Hit error while creating temp file: %s", err.Error())
			os.Remove(destfile.Name())
			ctx.JSON(apperrors.NewInternal().Status(), gin.H{
				//fmt.Fprintf(w, "Error occured during upload"),
				"error": err,
			})
			return
		}
		fileMetadata.TitleHash = hash.Original

		// defer the removal of the tempfile as well, something can be done
		// with it before the function is over (as long as you have the filehandle)
		//defer os.Remove(tempfile.Name())
		defer destfile.Close()
		//defer tempfile.Close()

		// here the temporary filename is known
		//log.Printf("Temporary filename: %s\n", tempfile.Name())
		log.Printf("Filename: %s\n", destfile.Name())

		for !uploaded {
			if n, err = part.Read(chunk); err != nil {
				if err != io.EOF {
					log.Printf("Hit error while reading chunk: %s", err.Error())
					os.Remove(destfile.Name())
					ctx.JSON(apperrors.NewInternal().Status(), gin.H{
						//fmt.Fprintf(w, "Error occured during upload"),
						"error": err,
					})
					return
				}
				uploaded = true
			}
			// if n, err = tempfile.Write(chunk[:n]); err != nil {}
			if n, err = destfile.Write(chunk[:n]); err != nil {
				log.Printf("Hit error while writing chunk: %s", err.Error())
				os.Remove(destfile.Name())
				ctx.JSON(apperrors.NewInternal().Status(), gin.H{
					//fmt.Fprintf(w, "Error occured during upload"),
					"error": err,
				})
				return
			}
			filesize += n
		}
		log.Printf("Uploaded filesize: %d bytes", filesize)

		// once uploaded something can be done with the file, the last defer
		// statement will remove the file after the function returns so any
		// errors during upload won't hit this, but at least the tempfile is
		// cleaned up
		log.Print(fileMetadata)

		//fileMetadata.Title = destfile.Name()

		// fileData := &models.File{
		// 	Title:    destfile.Name(),
		// 	AuthorID: authHeader.IDCode,
		// }
		log.Print("Hello fileData: ", fileMetadata)

		ctxRequest := ctx.Request.Context()

		err := handler.FilesService.UploadFile(ctxRequest, fileMetadata)
		if err != nil {
			log.Println("Failed to created Document record.")
			os.Remove(destfile.Name())
			ctx.JSON(apperrors.Status(err), gin.H{
				"error": err,
			})
			apperrors.NewInternal()
		}
	}
}

func (handler *handler) DeleteFile(ctx *gin.Context) {
	authHeader := ctx.MustGet("account").(*auth.JWTToken)
	accountID := uuid.MustParse(string(authHeader.IDCode.String()))
	log.Print(authHeader.IDCode)

	id := ctx.Params.ByName("id")
	filesID := uuid.MustParse(id)
	log.Print(filesID)
	//var request docRequest

	err := handler.FilesService.DeleteFile(ctx, filesID, accountID)
	if err != nil {
		log.Printf("error with deleting file %v", filesID)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":    "was unable to delete document",
			"document": filesID,
		})
		return
	}
	ctx.Writer.WriteHeader(204)
	//result, err := handler.
}

func (handler *handler) GetFile(ctx *gin.Context) {

}

func (handler *handler) GetFilesByDocID(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	documentID := uuid.MustParse(id)

	result, err := handler.FilesService.GetFileByDocID(ctx, documentID)
	if err != nil {
		log.Print("error in GetFileByDocID. Error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"erorr": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"files": result,
	})
}

func (handler *handler) UpdateFile(ctx *gin.Context) {
	authHeader := ctx.MustGet("account").(*auth.JWTToken)
	accountID := uuid.MustParse(string(authHeader.IDCode.String()))
	log.Print(authHeader.IDCode)

	id := ctx.Params.ByName("id")
	documentID := uuid.MustParse(id)
	log.Print(documentID)
	//var request docRequest

	err := handler.FilesService.UpdateFiles(ctx, documentID, accountID)
	if err != nil {
		log.Printf("error with deleting file %v", documentID)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":    "was unable to delete document",
			"document": documentID,
		})
		return
	}
	//result, err := handler.
}
