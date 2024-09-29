package gopherfiber

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

// ServerConfig holds the configuration for setting up the server.
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

// Server interface defines the behavior of a Fiber server.
type Server interface {
	Start() error
	GracefulShutdown()
	GetRouter() *fiber.App
}

// ServerSetup interface defines the setup methods for the Fiber server.
type ServerSetup interface {
	SetUpRouter(config ServerConfig) *fiber.App
	SetUpTLS(config ServerConfig) (*tls.Config, error)
	SetUpCORS(app *fiber.App, config ServerConfig)
}

// ServerSetupImpl is the concrete implementation of ServerSetup.
type ServerSetupImpl struct{}

// SetUpRouter sets up a Fiber server with static files and template paths.
func (s *ServerSetupImpl) SetUpRouter(config ServerConfig) *fiber.App {
	// Initialize the HTML template engine
	engine := html.New(config.TemplatePath, ".html") // Ensure config.TemplatePath points to your templates directory

	// Create a new Fiber app with the HTML template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files from the configured path
	app.Static("/static", config.StaticPath)

	return app
}

// SetUpTLS configures TLS (HTTPS) if enabled.
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
func (s *ServerSetupImpl) SetUpCORS(app *fiber.App, config ServerConfig) {
	if config.UseCORS {
		app.Use(cors.New(cors.Config{
			AllowOrigins: config.CORSConfig.AllowOrigins,
			AllowMethods: config.CORSConfig.AllowMethods,
		}))
		log.Printf("CORS configured with settings: %+v", config.CORSConfig)
	}
}

// FiberServer wraps the Fiber app and provides modular setup and shutdown.
type FiberServer struct {
	app         *fiber.App
	tlsConfig   *tls.Config
	serverSetup ServerSetup
	config      ServerConfig
}

// NewFiberServer creates a new FiberServer instance with injected dependencies.
func NewFiberServer(setup ServerSetup, config ServerConfig) Server {
	app := setup.SetUpRouter(config)
	setup.SetUpCORS(app, config)

	tlsConfig, err := setup.SetUpTLS(config)
	if err != nil {
		log.Fatalf("Error setting up TLS: %v", err)
	}

	return &FiberServer{
		app:         app,
		tlsConfig:   tlsConfig,
		serverSetup: setup,
		config:      config,
	}
}

// Start starts the Fiber server, with or without TLS.
func (fs *FiberServer) Start() error {
	addr := fmt.Sprintf(":%d", fs.config.Port)
	if fs.config.UseTLS {
		log.Printf("Starting server on port %d with TLS", fs.config.Port)
		go func() {
			if err := fs.app.ListenTLS(addr, fs.config.TLSCertFile, fs.config.TLSKeyFile); err != nil {
				log.Printf("ListenTLS error: %v", err)
			}
		}()
	} else {
		log.Printf("Starting server on port %d without TLS", fs.config.Port)
		go func() {
			if err := fs.app.Listen(addr); err != nil {
				log.Printf("Listen error: %v", err)
			}
		}()
	}

	return nil
}

// GetRouter returns the Fiber app instance.
func (fs *FiberServer) GetRouter() *fiber.App {
	return fs.app
}

func (fs *FiberServer) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	if err := fs.app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown successfully")
}

// Example usage:
//
//	func main() {
//	    config := gopherfiber.ServerConfig{
//	        Port:         8080,
//	        StaticPath:   "./static",
//	        TemplatePath: "./templates",
//	        UseTLS:       true,
//	        TLSCertFile:  "./certs/server.crt",
//	        TLSKeyFile:   "./certs/server.key",
//	        UseCORS:      true,
//	        CORSConfig: cors.Config{
//	            AllowOrigins: "https://example.com",
//	            AllowMethods: "GET,POST",
//	        },
//	    }
//
//	    server := gopherfiber.NewFiberServer(&gopherfiber.ServerSetupImpl{}, config)
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
