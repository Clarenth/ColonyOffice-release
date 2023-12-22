package main

import (
	"backend/config"
	// "backend/internal/handlers"
	// "backend/logger"

	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
	//_ "github.com/lib/pq"
)

const versionNumber = "0.0.1"

func main() {
	//gin.SetMode(gin.TestMode)
	// in here will be the server startup logic / spinup
	// spin up the server with the port, database connection, env file
	// logging will also be imported here to moinitor events, traffic

	log.Println("--->Loading config...")
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error: unable to complete configuration loading due to ", err)
	}
	log.Print("--->Config loaded successfully")

	log.Print("--->Injecting layers...")
	routerInjection, err := injection(config)
	if err != nil {
		log.Fatal("Error: Unable to complete data source injection due to ", err)
	}
	log.Print("--->Layers injected successfully")

	log.Println("Server port:", config.Port)
	gin.SetMode("debug")
	log.Println("Postgres connection:", config.DB)
	log.Println("Redis connection:", config.Redis)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: routerInjection,
	}

	// Initialize the server in a goroutine so that it won't block
	// the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// pprof profile and analysis
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Listen for the interrup signal
	<-ctx.Done()

	// Restore the default behavaiour on the interrupt signal and notify user of shutdown
	stop()
	log.Println("Graceful Shutdown initiated, press Ctrl+C again to force Shutdown")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown data storage connections
	if err := config.CloseDataStorageConnections(); err != nil {
		log.Fatalf("Possible error or Graceful Shutdown initiated. Closing data storage connections %v\n", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown due to: ", err)
	}
	log.Println("Server exiting...")
}
