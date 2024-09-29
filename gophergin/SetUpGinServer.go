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

// Server defines the behavior of a Gin server
type Server interface {
	Start() error
	GracefulShutdown()
	GetRouter() *gin.Engine
}

// ServerSetup defines the behavior for setting up a Gin server
type ServerSetup interface {
	SetUpRouter(config ServerConfig) *gin.Engine
	SetUpTLS(config ServerConfig) (*tls.Config, error)
	SetUpCORS(router *gin.Engine, config ServerConfig)
}

// ServerSetupImpl is the concrete implementation of ServerSetup
type ServerSetupImpl struct{}

// SetUpRouter sets up a Gin server with static and template paths
func (s *ServerSetupImpl) SetUpRouter(config ServerConfig) *gin.Engine {
	router := gin.Default()

	// Serve static files
	router.Static("/static", config.StaticPath)

	// Load HTML templates
	router.LoadHTMLGlob(config.TemplatePath)

	return router
}

// SetUpTLS configures the server for TLS if enabled
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

// SetUpCORS configures and applies CORS middleware if enabled
func (s *ServerSetupImpl) SetUpCORS(router *gin.Engine, config ServerConfig) {
	if config.UseCORS {
		router.Use(cors.New(config.CORSConfig))
		log.Printf("CORS configured with settings: %+v", config.CORSConfig)
	}
}

// GinServer is the modular implementation of the Server interface
type GinServer struct {
	router      *gin.Engine
	server      *http.Server
	serverSetup ServerSetup
	config      ServerConfig
}

// NewGinServer creates a new GinServer with injected ServerSetup
func NewGinServer(setup ServerSetup, config ServerConfig) Server {
	router := setup.SetUpRouter(config)
	setup.SetUpCORS(router, config)

	// Create the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	// Set TLS if enabled
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

// Start starts the Gin server with or without TLS
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

// GetRouter returns the gin.Engine instance
func (gs *GinServer) GetRouter() *gin.Engine {
	return gs.router
}

// GracefulShutdown handles graceful shutdown of the server
func (gs *GinServer) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gs.server.Shutdown(ctxShutDown); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server shutdown successfully")
}

// func GinServer() {
//     // Other initialization logic

//     // Create the server
//     serverConfig := gophergin.ServerConfig{
//         Port:         configs.ServerPort,
//         StaticPath:   "./static",
//         TemplatePath: "./templates",
//         UseTLS:       configs.UseTLS,
//         TLSCertFile:  configs.TLSCertFile,
//         TLSKeyFile:   configs.TLSKeyFile,
//         UseCORS:      configs.UseCORS,
//         CORSConfig:   cors.Config{
//             AllowOrigins: []string{"https://example.com"},
//         },
//     }

//     server := gophergin.NewGinServer(&gophergin.ServerSetupImpl{}, serverConfig)

//     // Get the router from the server
//     router := server.GetRouter()

//     // Add middleware
//     router.Use(middlewares.RequestIDGinMiddleware())

//     // Set up routes
//     routes.SetupSuperUserGinRoutes(router, superUserHandler, tokenManager)

//     // Start the server
//     err := server.Start()
//     if err != nil {
//         log.Fatalf("Failed to start server: %v", err)
//     }

//     // Gracefully shut down the server
//     server.GracefulShutdown()
// }
