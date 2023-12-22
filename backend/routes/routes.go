package routes

import (
	//"backend/internal/helpers/apperrors"
	//"backend/internal/middleware"
	//"backend/internal/models"
	mw "backend/routes/middleware"
	"time"

	auth_models "backend/modules/auth/models"
	documents_models "backend/modules/documents/models"
	files_models "backend/modules/files/models"

	"github.com/gin-gonic/gin"
)

type RoutesConfig struct {
	Router            *gin.Engine
	BaseURL           string
	TimeoutDuration   time.Duration
	DocumentDirectory string
	AuthHandler       auth_models.AuthHandlers
	TokenHandler      auth_models.TokenHandlers
	TokenService      auth_models.TokenService
	DocsHandler       documents_models.DocumentHandlers
	FilesHandler      files_models.FilesHandlers
}

/******Routes******/

// Routes initializes the endpoint routes and the handler with the required injected services.
// Does not return as it use a reference to the gin-gonic Engine.
func Routes(routes *RoutesConfig) {
	// Create a handler with injected services (later)
	authGroup := routes.Router.Group("/auth") // middleware.Timeout(config.TimeoutDuration, apperrors.NewServiceUnavailable())
	{
		if gin.Mode() != gin.TestMode {
			//authGroup.Use(middleware.Timeout(config.TimeoutDuration, apperrors.NewServiceUnavailable()))
			authGroup.GET("/account", mw.AuthAccount(routes.TokenService), routes.AuthHandler.CurrentAccount)
			authGroup.POST("/signout", mw.AuthAccount(routes.TokenService), routes.AuthHandler.Signout)
			authGroup.POST("/tokens", routes.TokenHandler.Tokens)
			authGroup.PATCH("/update", mw.AuthAccount(routes.TokenService), routes.AuthHandler.UpdateAccount)
			authGroup.DELETE("/delete/:id", mw.AuthAccount(routes.TokenService), routes.AuthHandler.DeleteAccount)
		} else {
			authGroup.GET("/account", routes.AuthHandler.CurrentAccount)
			authGroup.POST("/signout", routes.AuthHandler.Signout)

		}
		authGroup.POST("/signin", routes.AuthHandler.Signin)
		authGroup.POST("/signup", routes.AuthHandler.Signup)
	}

	apiGroup := routes.Router.Group(routes.BaseURL)
	{
		v1 := apiGroup.Group("/v1")
		{
			docsGroup := v1.Group("/docs") //.Use(mw.AuthAccount(config.TokenService))
			{
				if gin.Mode() != gin.TestMode {
					docsGroup.GET("/", routes.DocsHandler.GetAllDocuments)
					docsGroup.GET("?page=:index", routes.DocsHandler.GetDocsByPagination)                          // ?page=:index&count=:count
					docsGroup.POST("/add", mw.AuthAccount(routes.TokenService), routes.DocsHandler.CreateDocument) // mw.AuthAccount(config.TokenService)
					//docsGroup.POST("/upload", handler.UploadFiles)                                      // mw.AuthAccount(config.TokenService)
					docsGroup.GET("/:id", mw.AuthAccount(routes.TokenService), routes.DocsHandler.GetDocumentByID)
					docsGroup.POST("/search", mw.AuthAccount(routes.TokenService), routes.DocsHandler.SearchDocs)
					docsGroup.DELETE("/delete/:id", mw.AuthAccount(routes.TokenService), routes.DocsHandler.DeleteDoc)
					docsGroup.PATCH("/update/:id", mw.AuthAccount(routes.TokenService), routes.DocsHandler.UpdateDoc)
				} else {
					docsGroup.GET("/", routes.DocsHandler.GetAllDocuments)
					docsGroup.POST("/add", routes.DocsHandler.CreateDocument) // mw.AuthAccount(config.TokenService)
					docsGroup.POST("/search", routes.DocsHandler.SearchDocs)
				}
			}
			filesGroup := v1.Group("/files") //.Use(mw.AuthAccount(config.TokenService))
			{
				if gin.Mode() != gin.TestMode {
					filesGroup.POST("/upload", mw.AuthAccount(routes.TokenService), routes.FilesHandler.UploadFiles)
					filesGroup.DELETE("/delete/:id", mw.AuthAccount(routes.TokenService), routes.FilesHandler.DeleteFile)
					filesGroup.GET("/:id", mw.AuthAccount(routes.TokenService), routes.FilesHandler.GetFile)
					filesGroup.GET("/doc-id/:id", mw.AuthAccount(routes.TokenService), routes.FilesHandler.GetFilesByDocID)
					filesGroup.PATCH("/update/:id", routes.FilesHandler.UpdateFile)
				} else {
					filesGroup.POST("/upload", routes.FilesHandler.UploadFiles)
					filesGroup.POST("/delete/:id", routes.FilesHandler.UploadFiles)
					filesGroup.PATCH("/update/:id", routes.FilesHandler.UpdateFile)

				}
			}
		}

	}

}

/******MiddleWare******/
