package gophergin

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ServerConfig holds the configuration for setting up the server
type ServerConfig struct {
	Port         int
	StaticPath   string
	TemplatePath string
	UseTLS       bool
	TLSCertFile  string
	TLSKeyFile   string
	UseCORS      bool
	CORSConfig   cors.Config
}

// ServerSetup defines the interface for setting up a Gin server
type ServerSetup interface {
	SetUpServer(config ServerConfig) (*gin.Engine, error)
}

// CorsServerSetup is the struct for setting up the server with CORS
type CorsServerSetup struct{}

// SetUpServer sets up a Gin server with CORS middleware.
//
// This function initializes a Gin server with the provided static and template paths,
// and applies CORS middleware if configured.
//
// Example usage:
//
//	serverConfig := ServerConfig{
//	    StaticPath:   "./static",
//	    TemplatePath: "./templates",
//	    UseCORS:      true,
//	    CORSConfig: cors.Config{
//	        AllowOrigins: []string{"https://example.com"},
//	    },
//	}
//	router, err := (&CorsServerSetup{}).SetUpServer(serverConfig)
//	if err != nil {
//	    log.Fatalf("Failed to set up server: %v", err)
//	}
func (c *CorsServerSetup) SetUpServer(config ServerConfig) (*gin.Engine, error) {
	router := gin.Default()

	// Serve static files
	router.Static("/static", config.StaticPath)

	// Load HTML templates
	router.LoadHTMLGlob(config.TemplatePath)

	// Configure and apply CORS middleware if enabled
	if config.UseCORS {
		router.Use(cors.New(config.CORSConfig))
		log.Printf("CORS configured with settings: %+v", config.CORSConfig)
	}

	return router, nil
}

// BasicServerSetup is the struct for setting up the server without CORS
type BasicServerSetup struct{}

// SetUpServer sets up a basic Gin server without CORS middleware.
//
// Example usage:
//
//	serverConfig := ServerConfig{
//	    StaticPath:   "./static",
//	    TemplatePath: "./templates",
//	}
//	router, err := (&BasicServerSetup{}).SetUpServer(serverConfig)
//	if err != nil {
//	    log.Fatalf("Failed to set up server: %v", err)
//	}
func (b *BasicServerSetup) SetUpServer(config ServerConfig) (*gin.Engine, error) {
	router := gin.Default()

	// Serve static files
	router.Static("/static", config.StaticPath)

	// Load HTML templates
	router.LoadHTMLGlob(config.TemplatePath)

	return router, nil
}

// StartGinServer starts the provided Gin server with or without TLS.
//
// This function starts a Gin server using the provided *gin.Engine and ServerConfig.
// If TLS is enabled, the server will use the provided certificate and key files for HTTPS.
//
// Example usage:
//
//	err := StartGinServer(router, serverConfig)
//	if err != nil {
//	    log.Fatalf("Failed to start server: %v", err)
//	}
func StartGinServer(router *gin.Engine, config ServerConfig) error {
	// Create the HTTP server with the provided router and port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	if config.UseTLS {
		// Load the TLS certificate and key
		cert, err := LoadTLSCertificate(config.TLSCertFile, config.TLSKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}

		// Configure the server for TLS
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		log.Printf("Gin server is running on port %d with TLS", config.Port)

		// Start the server with TLS
		go func() {
			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Printf("ListenAndServeTLS error: %v", err)
			}
		}()
	} else {
		log.Printf("Gin server is running on port %d without TLS", config.Port)

		// Start the server without TLS
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("ListenAndServe error: %v", err)
			}
		}()
	}

	return nil
}

// GracefulShutdown handles the graceful shutdown of the Gin server.
//
// This function listens for interrupt signals (Ctrl+C) and shuts down the server gracefully,
// allowing any in-flight requests to complete within a timeout period.
//
// Example usage:
//
//	go GracefulShutdown(server)
func GracefulShutdown(server *http.Server) {
	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Context with timeout for shutdown
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Printf("Server forced to shutdown: %v", err)

		// Retry mechanism for shutdown
		for retries := 0; retries < 3; retries++ {
			log.Printf("Retrying shutdown... attempt %d", retries+1)
			if err := server.Shutdown(ctxShutDown); err == nil {
				log.Println("Server shutdown successfully on retry")
				return
			}
		}
		log.Fatalf("Failed to shutdown server gracefully after retries: %v", err)
	}

	log.Println("Server shutdown successfully")
}

// LoadTLSCertificate loads the TLS certificate and private key.
//
// This function loads a TLS certificate and key from the specified files and returns
// a tls.Certificate object that can be used to configure HTTPS servers.
//
// Example usage:
//
//	cert, err := LoadTLSCertificate("/path/to/cert.crt", "/path/to/key.key")
//	if err != nil {
//	    log.Fatalf("Failed to load TLS certificate: %v", err)
//	}
func LoadTLSCertificate(certFile, keyFile string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load TLS certificate: %w", err)
	}
	return cert, nil
}

// func main() {
// 	// Define server configuration
// 	serverConfig := server.ServerConfig{
// 		Port:        8080,
// 		StaticPath:  "./static",
// 		TemplatePath: "./templates",
// 		UseTLS:      false,
// 		UseCORS:     true,
// 		CORSConfig: cors.Config{
// 			AllowOrigins: []string{"https://example.com"},
// 			AllowMethods: []string{"GET", "POST"},
// 		},
// 	}

// 	// Set up the server with CORS
// 	router, err := (&server.CorsServerSetup{}).SetUpServer(serverConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to set up server: %v", err)
// 	}

// 	// Start the server
// 	if err := server.StartGinServer(router, serverConfig); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}

// 	// Gracefully shutdown on interrupt
// 	server.GracefulShutdown(&http.Server{Addr: ":8080"})
// }
