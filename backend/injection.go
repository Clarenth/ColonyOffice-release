package main

import (
	"backend/config"
	"backend/logger"
	"backend/routes"

	// auth module
	auth_handlers "backend/modules/auth/handlers"
	auth_repository "backend/modules/auth/repository"
	auth_services "backend/modules/auth/services"

	// documents module
	documents_handlers "backend/modules/documents/handlers"
	documents_repository "backend/modules/documents/repository"
	documents_services "backend/modules/documents/services"

	// files module
	files_handlers "backend/modules/files/handlers"
	files_repository "backend/modules/files/repository"
	files_services "backend/modules/files/services"

	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/jackc/pgx/v5"
)

func injection(config *config.Configuration) (*gin.Engine, error) {
	// RSA Keys
	log.Println("Loading RSA keys")

	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	privateKeyValue, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key file: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyValue)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	publicKeyFile := os.Getenv("PUBLIC_KEY_FILE")
	publicKeyValue, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key file: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyValue)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// JWT tokens
	id_token_exp := os.Getenv("ID_TOKEN_EXP")
	refresh_token_exp := os.Getenv("REFRESH_TOKEN_EXP")
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	idTokenExpiration, err := strconv.ParseInt(id_token_exp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse JWT token expiration ENV to int: %w", err)
	}
	refreshTokenExpiration, err := strconv.ParseInt(refresh_token_exp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse JWT refresh expiration ENV to int: %w", err)
	}
	log.Print("idToken expiration time: ", idTokenExpiration)
	log.Print("refreshToken expiration time: ", refreshTokenExpiration)

	// Repository layer
	accountRepository := auth_repository.NewAccountRepository(config.DB)
	tokenRepository := auth_repository.NewRedisTokenRepository(config.Redis)
	documentRepository := documents_repository.NewDocumentRepository(config.DB)
	filesRepository := files_repository.NewFilesRepository(config.DB)

	// Services
	accountServices := auth_services.NewEmployeeService(&auth_services.ConfigAccountService{
		EmployeeRepository: accountRepository,
	})
	tokenServices := auth_services.NewTokenService(&auth_services.ConfigTokenService{
		TokenRepository:            tokenRepository,
		PrivateKey:                 privateKey,
		PublicKey:                  publicKey,
		RefreshSecretKey:           refreshSecretKey,
		IDTokenExpirationSecs:      idTokenExpiration,
		RefreshTokenExpirationSecs: refreshTokenExpiration,
	})
	documentServices := documents_services.NewDocumentService(&documents_services.ConfigDocumentService{
		DocumentDirectory:  config.FilesDir,
		DocumentRepository: documentRepository,
	})
	filesServices := files_services.NewFilesService(&files_services.ConfigFilesService{
		FilesRepository: filesRepository,
	})

	// handlers
	authHandlers := auth_handlers.NewAuthHandlers(&auth_handlers.HandlerConfig{
		EmployeeService: accountServices,
		TokenService:    tokenServices,
	})
	documentHandlers := documents_handlers.NewDocumentHandlers(&documents_handlers.HandlerConfig{
		DocumentService: documentServices,
	})
	filesHandlers := files_handlers.NewFilesHandlers(&files_handlers.HandlerConfig{
		DocumentDirectory: config.FilesDir,
		FilesService:      filesServices,
	})

	log.Println("Injecting data sources for layers")

	// Server setup
	baseURL := os.Getenv("BASE_URL")

	handler_timeout := os.Getenv("HANDLER_TIMEOUT")
	handlerTimeout, err := strconv.ParseInt(handler_timeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("error: could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	// CORS configuration must be done before router initialization
	corsConfig := cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "https://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Accept", "X-CSRF-Token", "Authorization", "Content-Type", "access-control-allow-origin", "ACCESS-CONTROL-ALLOW-ORIGIN", "Access-Control-Allow-Headers", "Origin", "origin"},
	}

	router := gin.New()
	router.MaxMultipartMemory = 2 << 20 // sets the memory allocated per each parse (2MB, default is 32MB), not how large an uploaded file can be
	router.Use(cors.New(corsConfig), gin.Recovery(), logger.CustomLogger(config.LoggerDir))

	routes.Routes(&routes.RoutesConfig{
		Router:            router,
		BaseURL:           baseURL,
		TimeoutDuration:   time.Duration(time.Duration(handlerTimeout) * time.Second),
		DocumentDirectory: config.FilesDir,
		AuthHandler:       authHandlers,
		TokenHandler:      authHandlers,
		TokenService:      tokenServices,
		DocsHandler:       documentHandlers,
		FilesHandler:      filesHandlers,
	})

	return router, nil
}

// Use in the server graceful shutdown
// Close all remote Data Storage connections (Postgres, Redis, Cloud, etc.)
// func (config *config.Configuration) CloseDataStorageConnections() error {
// 	if err := config.DB.Close(); err != nil {
// 		return fmt.Errorf("error, closing Postgres database: %w", err)
// 	}

// 	//Redis
// 	if err := config.Redis.Close(); err != nil {
// 		return fmt.Errorf("error, closing Redis database:%w", err)
// 	}

// 	//Cloud Storage(?)

// 	//CDN for files

// 	return nil
// }
