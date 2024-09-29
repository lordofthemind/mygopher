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

// ServerConfig holds the configuration for setting up the server.
//
// Fields:
// - Port: Port number to run the server on.
// - StaticPath: Path to serve static files from.
// - TemplatePath: Path to HTML templates for rendering.
// - UseTLS: Enable TLS (HTTPS) if true.
// - TLSCertFile: Path to the TLS certificate file (required if UseTLS is true).
// - TLSKeyFile: Path to the TLS key file (required if UseTLS is true).
// - UseCORS: Enable CORS (Cross-Origin Resource Sharing) if true.
// - CORSConfig: Configures allowed origins, headers, and methods for CORS.
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

// Server interface defines the behavior of a Gin server.
//
// Methods:
// - Start: Starts the server (optionally with TLS).
// - GracefulShutdown: Gracefully shuts down the server when interrupted.
// - GetRouter: Returns the underlying gin.Engine for additional route setup.
type Server interface {
	Start() error
	GracefulShutdown()
	GetRouter() *gin.Engine
}

// ServerSetup defines the behavior for setting up a Gin server.
//
// Methods:
// - SetUpRouter: Configures and returns a new Gin engine with static file and template paths.
// - SetUpTLS: Configures TLS settings if required (returns a tls.Config instance).
// - SetUpCORS: Applies CORS middleware to the Gin engine if enabled.
type ServerSetup interface {
	SetUpRouter(config ServerConfig) *gin.Engine
	SetUpTLS(config ServerConfig) (*tls.Config, error)
	SetUpCORS(router *gin.Engine, config ServerConfig)
}

// ServerSetupImpl is the concrete implementation of ServerSetup.
type ServerSetupImpl struct{}

// SetUpRouter sets up a Gin server with static file and template paths.
//
// Parameters:
// - config: The server configuration for static files and template paths.
//
// Returns:
// - *gin.Engine: A configured Gin engine.
func (s *ServerSetupImpl) SetUpRouter(config ServerConfig) *gin.Engine {
	router := gin.Default()

	// Serve static files from the configured path.
	router.Static("/static", config.StaticPath)

	// Load HTML templates from the configured path.
	router.LoadHTMLGlob(config.TemplatePath)

	return router
}

// SetUpTLS configures the server for TLS (HTTPS) if enabled.
//
// Parameters:
// - config: The server configuration containing TLS settings.
//
// Returns:
// - *tls.Config: TLS configuration if enabled, or nil if not.
// - error: An error if TLS certificates cannot be loaded.
func (s *ServerSetupImpl) SetUpTLS(config ServerConfig) (*tls.Config, error) {
	if !config.UseTLS {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(config.TLSCertFile, config.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return tlsConfig, nil
}

// SetUpCORS configures and applies CORS middleware if enabled.
//
// Parameters:
// - router: The Gin engine to apply the middleware to.
// - config: The server configuration that contains CORS settings.
func (s *ServerSetupImpl) SetUpCORS(router *gin.Engine, config ServerConfig) {
	if config.UseCORS {
		router.Use(cors.New(config.CORSConfig))
		log.Printf("CORS configured with settings: %+v", config.CORSConfig)
	}
}

// GinServer is the modular implementation of the Server interface.
// It wraps around Gin's HTTP server and provides modular setup and shutdown.
type GinServer struct {
	router      *gin.Engine
	server      *http.Server
	serverSetup ServerSetup
	config      ServerConfig
}

// NewGinServer creates a new GinServer instance with injected dependencies.
//
// Parameters:
// - setup: A ServerSetup implementation for initializing the server.
// - config: The ServerConfig structure for server configuration.
//
// Returns:
// - Server: A configured Gin server ready to start.
func NewGinServer(setup ServerSetup, config ServerConfig) Server {
	router := setup.SetUpRouter(config)
	setup.SetUpCORS(router, config)

	// Create the HTTP server instance.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	// Set up TLS if enabled.
	tlsConfig, err := setup.SetUpTLS(config)
	if err != nil {
		log.Fatalf("Error setting up TLS: %v", err)
	}
	server.TLSConfig = tlsConfig

	return &GinServer{
		router:      router,
		server:      server,
		serverSetup: setup,
		config:      config,
	}
}

// Start starts the Gin server, either with or without TLS.
//
// Returns:
// - error: Any error encountered while starting the server.
func (gs *GinServer) Start() error {
	if gs.config.UseTLS {
		log.Printf("Starting server on port %d with TLS", gs.config.Port)
		go func() {
			if err := gs.server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Printf("ListenAndServeTLS error: %v", err)
			}
		}()
	} else {
		log.Printf("Starting server on port %d without TLS", gs.config.Port)
		go func() {
			if err := gs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("ListenAndServe error: %v", err)
			}
		}()
	}

	return nil
}

// GetRouter returns the gin.Engine instance.
//
// Returns:
// - *gin.Engine: The underlying Gin engine for the server.
func (gs *GinServer) GetRouter() *gin.Engine {
	return gs.router
}

// GracefulShutdown gracefully shuts down the server when interrupted.
//
// This method handles system interrupts (e.g., Ctrl+C) and shuts down the server
// gracefully, allowing for ongoing requests to finish within a 5-second timeout.
func (gs *GinServer) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Allow up to 5 seconds for graceful shutdown.
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gs.server.Shutdown(ctxShutDown); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server shutdown successfully")
}

// Example usage:
//
//	func main() {
//	    config := gophergin.ServerConfig{
//	        Port:         8080,
//	        StaticPath:   "./static",
//	        TemplatePath: "./templates/*.html",
//	        UseTLS:       true,
//	        TLSCertFile:  "./certs/server.crt",
//	        TLSKeyFile:   "./certs/server.key",
//	        UseCORS:      true,
//	        CORSConfig: cors.Config{
//	            AllowOrigins: []string{"https://example.com"},
//	            AllowMethods: []string{"GET", "POST"},
//	        },
//	    }
//
//	    server := gophergin.NewGinServer(&gophergin.ServerSetupImpl{}, config)
//
//	    // Start the server
//	    err := server.Start()
//	    if err != nil {
//	        log.Fatalf("Failed to start server: %v", err)
//	    }
//
//	    // Gracefully shut down the server on interrupt
//	    server.GracefulShutdown()
//	}
