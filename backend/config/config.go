package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	//_ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/peterbourgon/ff"
	"github.com/redis/go-redis/v9"
)

/*
	// For enviroment
	[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
	- using env:   export GIN_MODE=release
	- using code:  gin.SetMode(gin.ReleaseMode)
*/

// Configuration struct fields are used for configuring the server
type Configuration struct {
	Port       string
	Enviroment string
	DB         *sqlx.DB
	Redis      *redis.Client
	FilesDir   string
	LoggerDir  string
}

// LoadConfig initializes a new Configuration singleton instance and loads the .env file
// along with loading the databases, directories.
func LoadConfig() (*Configuration, error) {
	envLoadError := godotenv.Load(".env.dev")
	if envLoadError != nil {
		log.Fatal("Error loading env file: ", envLoadError)
	}

	port := os.Getenv("PORT")
	enviroment := os.Getenv("ENVIROMENT")
	pgUser := os.Getenv("DB_USERNAME")
	pgPassword := os.Getenv("DB_PASSWORD")
	pgHost := os.Getenv("DB_HOST")
	pgPort := os.Getenv("DB_PORT")
	pgName := os.Getenv("DB_NAME")
	pgSSL := os.Getenv("DB_SSL")
	pgConnectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", pgUser, pgPassword, pgHost, pgPort, pgName, pgSSL)
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")

	// Postgres verify connection
	log.Printf("Config: Connecting to Postgres")
	postgres, err := sqlx.Connect("pgx", pgConnectionString)
	postgres.SetMaxIdleConns(5)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Redis verify connection
	log.Printf("Connecting to redis")
	redisDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redis_host, redis_port),
		Password: "",
		DB:       0,
	})

	// Files directory
	filesDir := os.Getenv("FILES_DIRECTORY")
	if _, err := os.Stat(filesDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filesDir, 0700)
		if err != nil {
			log.Println(err)
		}
	}
	log.Printf("Document files saved to directory: %v", filesDir)

	// Logger Files directory
	loggerDir := os.Getenv("LOGS_DIRECTORY")
	if _, err := os.Stat(loggerDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(loggerDir, 0700)
		if err != nil {
			log.Println(err)
		}
	}
	log.Printf("Logger files saved to directory: %v", loggerDir)

	return &Configuration{
		Port:       port,
		Enviroment: enviroment,
		DB:         postgres,
		Redis:      redisDB,
		FilesDir:   filesDir,
		LoggerDir:  loggerDir,
	}, nil
}

// Use in the server graceful shutdown to close all remote Data Storage connections (Postgres, Redis, Cloud, etc.)
func (config *Configuration) CloseDataStorageConnections() error {
	if err := config.DB.Close(); err != nil {
		return fmt.Errorf("error, closing Postgres database: %w", err)
	}

	//Redis
	if err := config.Redis.Close(); err != nil {
		return fmt.Errorf("error, closing Redis database:%w", err)
	}

	//Cloud Storage(?)

	//CDN for files(?)

	return nil
}

func SetFlagsNoEnv() {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)

	port := fs.String("port", os.Getenv("PORT"), "Server port to listn on")
	enviroment := fs.String("env", os.Getenv("ENVIROMENT"), "Application enviroment (development | production)")
	dbUser := fs.String("pgUser", os.Getenv("DB_USERNAME"), "Database username")
	dbPassword := fs.String("pgPassword", os.Getenv("DB_Password"), "Database password")
	dbHost := fs.String("pgHost", os.Getenv("DB_Host"), "Database host name")
	dbPort := fs.String("pgPort", os.Getenv("DB_PORT"), "Database port number")
	dbName := fs.String("pgName", os.Getenv("DB_NAME"), "Database server name")
	dbSSL := fs.String("pgSSLSetting", os.Getenv("DB_SSL"), "Database SSL setting")

	pgConnectionString := fmt.Sprintf("user=%p password=%p host=%p port=%p name=%p ssl=%p", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSL)
	//fs.StringVar(&Config.DB_Connection, "postgres DB", pgConnectionString, "PostgreSQL connection")
	//fs.StringVar(&Config.jwt.secret, "jwt", "jwt-secret", "jwt authorization string")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("port %d, debug %v\n, postgres: %v", &port, &enviroment, &pgConnectionString)
	//fmt.Printf("port %s, debug %v\n", *port, *debug)

	log.Printf("Connecting to Postgres")
	db, err := sqlx.Open("postgres", pgConnectionString)
	if err != nil {
		fmt.Errorf("error opening db: %w", err)
	}

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		fmt.Errorf("error connecting to db: %w", err)
	}
}

/*
func LoadConfig() {
	file, err := ff.WithConfigFile(".env.dev")
	if err != nil {
		log.Fatal("Error, file not found.")
		return err
	}
}
*/
