package handlers

import (
	auth "backend/modules/auth/models"
	"backend/modules/documents/helpers/apperrors"
	"backend/modules/documents/models"
	"strconv"

	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler holds the methods expected by the service layer for handling routing
type handler struct {
	DocumentService models.DocumentService
}

// HandlerConfig holds the server configuration need by by the handler
type HandlerConfig struct {
	DocumentService models.DocumentService
}

func NewDocumentHandlers(ctx *HandlerConfig) models.DocumentHandlers {
	return &handler{
		DocumentService: ctx.DocumentService,
	}
}

func (handler *handler) FormFileData(ctx *gin.Context) {
	panic("No longer in use")
	// authHeader := ctx.MustGet("account").(*auth.JWTToken)
	// var request docFormRequest

	// // Bind incoming JSON to a struct and check for validation errors
	// if ok := bindData(ctx, &request); !ok {
	// 	//log.Print(request)
	// 	return
	// }

	// newDocumentData := &models.DocumentWebForm{
	// 	DocumentTitle:       request.Title,
	// 	AuthorName:          request.Author,
	// 	AuthorID:            authHeader.IDCode,
	// 	Description:         request.Description,
	// 	DocumentContent:     request.Content,
	// 	ClassificationLevel: request.SecurityAccessLevel,
	// 	/*ClassificationLevel: &models.SecurityAccessLevel{
	// 		ClassificationLevel: request.SecurityAccessLevel.ClassificationLevel,
	// 	},*/
	// 	Language: ctx.Request.Header.Get("Accept-Language"),
	// }
	// log.Print(newDocumentData)
	// //ctx.DefaultPostForm()
	// ctxRequest := ctx.Request.Context()
	// err := handler.DocumentService.WriteFormToFile(ctxRequest, newDocumentData)
	// if err != nil {
	// 	log.Println("Failed to created Document record.")
	// 	ctx.JSON(apperrors.Status(err), gin.H{
	// 		"error": err,
	// 	})
	// 	apperrors.NewInternal()
	// }

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"result": "file saved successfully",
	// })

}

func (handler *handler) CreateDocument(ctx *gin.Context) {
	authHeader := ctx.MustGet("account").(*auth.JWTToken)
	var request docRequest

	// Bind incoming JSON to a struct and check for validation errors
	if ok := bindData(ctx, &request); !ok {
		//log.Print(request)
		return
	}

	documentData := &models.Document{
		DocumentTitle: request.Title,
		//AuthorName:          request.Author,
		AuthorID:            authHeader.IDCode,
		Description:         request.Description,
		Language:            request.Language,
		ClassificationLevel: request.SecurityAccessLevel,
	}

	ctxRequest := ctx.Request.Context()
	documentResponse, err := handler.DocumentService.WriteDocumentForm(ctxRequest, documentData)
	if err != nil {
		log.Println("Failed to created Document record.")
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		apperrors.NewInternal()
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"document_code": documentResponse,
	})
}

func (handler *handler) DeleteDoc(ctx *gin.Context) {
	authHeader := ctx.MustGet("account").(*auth.JWTToken)
	accountID := uuid.MustParse(string(authHeader.IDCode.String()))
	log.Print(authHeader.IDCode)

	id := ctx.Params.ByName("id")
	documentID := uuid.MustParse(id)
	log.Print(documentID)

	err := handler.DocumentService.DeleteDocument(ctx, documentID, accountID)
	if err != nil {
		log.Print("Could not delete document")
	}

	response := fmt.Sprintf("Document %v was deleted", documentID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"result": response,
	})
}

func (handler *handler) UpdateDoc(ctx *gin.Context) {
	authHeader := ctx.MustGet("account").(*auth.JWTToken)
	accountID := uuid.MustParse(string(authHeader.IDCode.String()))
	log.Print(authHeader.IDCode)

	id := ctx.Params.ByName("id")
	documentID := uuid.MustParse(id)
	log.Print(documentID)
	var request docUpdateRequest

	if ok := bindData(ctx, &request); !ok {
		log.Print("error binding Document request to JSON.")
		return
	}

	docRequest := &models.Document{
		DocumentID:          request.DocumentID,
		DocumentTitle:       request.DocumentTitle,
		AuthorName:          request.AuthorName,
		AuthorID:            authHeader.IDCode,
		Description:         request.Description,
		ClassificationLevel: request.ClassificationLevel,
		Language:            request.Language,
	}

	updateResponse, err := handler.DocumentService.UpdateDocument(ctx, docRequest, accountID)
	if err != nil {
		log.Print("error with updating document service layer")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "there was an error with updating the document",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": updateResponse,
	})
}

func (handler *handler) GetAllDocuments(ctx *gin.Context) {
	documents, err := handler.DocumentService.GetAllDocumentRecords(ctx)
	if err != nil {
		ctx.JSON(apperrors.NewInternal().Status(), gin.H{
			"error": "could not get requested documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"documents": documents,
		//"result": "endpoint reach successfuly",
	})
	//return
}

func (handler *handler) GetDocumentByID(ctx *gin.Context) {
	// authHeader := ctx.MustGet("account").(*auth.JWTToken)
	// log.Print(authHeader.IDCode)

	id := ctx.Params.ByName("id")
	documentID := uuid.MustParse(id)
	log.Print(documentID)

	document, err := handler.DocumentService.GetDocumentByID(ctx, documentID.String())
	if err != nil {
		ctx.JSON(apperrors.NewInternal().Status(), gin.H{
			"error": "could not get requested document",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"document": document,
		//"result": "endpoint reach successfuly",
	})
	//return
}

func (handler *handler) GetDocsByPagination(ctx *gin.Context) {
	pageIndex, err := strconv.Atoi((ctx.Params.ByName("page")))
	if err != nil {
		log.Printf("error with getting URL param %s", pageIndex)
	}
	pageCount, err := strconv.Atoi((ctx.Params.ByName("page")))
	if err != nil {
		log.Printf("error with getting URL param %v", pageCount)
	}

	log.Print(pageIndex)
	log.Print(pageCount)

	result, err := handler.DocumentService.GetDocsByPagination(ctx, pageIndex, pageCount)
	if err != nil {
		log.Print(err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.JSON(200, gin.H{
		"documents": result,
	})

}

func (handler *handler) SearchDocs(ctx *gin.Context) {
	//panic("Func not implemented yet")
	//authHeader := ctx.MustGet("account").(*auth.JWTToken)
	var request docGetRequest_Single
	log.Printf("Hello %s request", request)

	// Bind incoming JSON to a struct and check for validation errors
	if ok := bindData(ctx, &request); !ok {
		log.Print(request)
		return
	}
	//log.Printf("Hello %s from GetOneFile", authHeader)
	log.Printf("Hello %s request", request)

	getDocsRequest := &models.DocumentSearchRequest{
		SearchTags:          request.SearchTags,
		Author:              request.Author,
		Colony:              request.Colony,
		DatesRangeStart:     request.DatesRangeStart,
		DatesRangeEnd:       request.DatesRangeEnd,
		SecurityAccessLevel: request.SecurityAccessLevel,
	}

	log.Print("Hello from handler: ", getDocsRequest)

	documents, err := handler.DocumentService.SearchDocsByTitle(ctx, getDocsRequest)
	if err != nil {
		ctx.JSON(apperrors.NewInternal().Status(), gin.H{
			"error": "could not get requested documents",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"documents": documents,
		//"result": "endpoint reach successfuly",
	})
	//return
}
