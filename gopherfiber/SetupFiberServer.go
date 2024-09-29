package gopherfiber

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// ServerConfig holds the configuration options for the server.
//
// Fields:
// - Port: The port number where the server will listen for incoming requests.
// - UseTLS: Set to true if you want to enable HTTPS using TLS.
// - TLSCertFile: Path to the TLS certificate file (required if UseTLS is true).
// - TLSKeyFile: Path to the TLS key file (required if UseTLS is true).
// - UseCORS: Set to true to enable Cross-Origin Resource Sharing (CORS).
// - CORSConfig: CORS configuration to allow specific origins and methods.
type ServerConfig struct {
	Port        int
	UseTLS      bool
	TLSCertFile string
	TLSKeyFile  string
	UseCORS     bool
	CORSConfig  cors.Config
}

// Server interface defines the behavior of a Fiber server.
//
// Methods:
// - Start: Starts the server (optionally with TLS).
// - GracefulShutdown: Gracefully shuts down the server on interrupt.
// - GetRouter: Returns the underlying fiber.App instance for adding routes.
type Server interface {
	Start() error
	GracefulShutdown()
	GetRouter() *fiber.App
}

// ServerSetup interface defines the setup methods for configuring a Fiber server.
//
// Methods:
// - SetUpRouter: Sets up the Fiber app without static files or templates.
// - SetUpTLS: Configures TLS settings for HTTPS if enabled.
// - SetUpCORS: Applies CORS middleware to the app if enabled.
type ServerSetup interface {
	SetUpRouter(config ServerConfig) *fiber.App
	SetUpTLS(config ServerConfig) (*tls.Config, error)
	SetUpCORS(app *fiber.App, config ServerConfig)
}

// ServerSetupImpl is the concrete implementation of the ServerSetup interface.
type ServerSetupImpl struct{}

// SetUpRouter sets up a Fiber app without static file serving or template rendering.
//
// Parameters:
// - config: The ServerConfig structure.
//
// Returns:
// - *fiber.App: The Fiber app instance.
func (s *ServerSetupImpl) SetUpRouter(config ServerConfig) *fiber.App {
	// Create a new Fiber app
	app := fiber.New()

	return app
}

// SetUpTLS configures TLS (HTTPS) settings if enabled.
//
// Parameters:
// - config: The ServerConfig containing the paths to the TLS certificate and key files.
//
// Returns:
// - *tls.Config: The TLS configuration if UseTLS is true, or nil if not.
// - error: An error if the certificate or key cannot be loaded.
func (s *ServerSetupImpl) SetUpTLS(config ServerConfig) (*tls.Config, error) {
	if !config.UseTLS {
		return nil, nil
	}

	// Load TLS certificate and key
	cert, err := tls.LoadX509KeyPair(config.TLSCertFile, config.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Return the TLS configuration
	return tlsConfig, nil
}

// SetUpCORS configures and applies CORS middleware to the app if enabled.
//
// Parameters:
// - app: The Fiber app to apply CORS middleware to.
// - config: The ServerConfig containing CORS configuration options.
func (s *ServerSetupImpl) SetUpCORS(app *fiber.App, config ServerConfig) {
	if config.UseCORS {
		// Apply CORS middleware with the provided configuration
		app.Use(cors.New(cors.Config{
			AllowOrigins: config.CORSConfig.AllowOrigins,
			AllowMethods: config.CORSConfig.AllowMethods,
		}))
		log.Printf("CORS configured with settings: %+v", config.CORSConfig)
	}
}

// FiberServer wraps the Fiber app and provides modular setup and graceful shutdown.
//
// Fields:
// - app: The Fiber app instance used for routing and middleware.
// - tlsConfig: Optional TLS configuration for serving HTTPS.
// - serverSetup: The ServerSetup instance used to configure the server.
// - config: The ServerConfig structure containing the server's configuration.
type FiberServer struct {
	app         *fiber.App
	tlsConfig   *tls.Config
	serverSetup ServerSetup
	config      ServerConfig
}

// NewFiberServer creates a new FiberServer instance with the provided configuration.
//
// Parameters:
// - setup: The ServerSetup implementation to configure the server.
// - config: The ServerConfig structure containing configuration details.
//
// Returns:
// - Server: A configured Fiber server instance ready to start.
func NewFiberServer(setup ServerSetup, config ServerConfig) Server {
	// Initialize the Fiber app without static file serving or templates
	app := setup.SetUpRouter(config)
	// Configure CORS if enabled
	setup.SetUpCORS(app, config)

	// Set up TLS if enabled
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

// Start starts the Fiber server with or without TLS, depending on the configuration.
//
// Returns:
// - error: Any error encountered during server startup.
func (fs *FiberServer) Start() error {
	addr := fmt.Sprintf(":%d", fs.config.Port)
	if fs.config.UseTLS {
		// Start the server with TLS
		log.Printf("Starting server on port %d with TLS", fs.config.Port)
		go func() {
			if err := fs.app.ListenTLS(addr, fs.config.TLSCertFile, fs.config.TLSKeyFile); err != nil {
				log.Printf("ListenTLS error: %v", err)
			}
		}()
	} else {
		// Start the server without TLS
		log.Printf("Starting server on port %d without TLS", fs.config.Port)
		go func() {
			if err := fs.app.Listen(addr); err != nil {
				log.Printf("Listen error: %v", err)
			}
		}()
	}

	return nil
}

// GetRouter returns the underlying Fiber app instance.
//
// Returns:
// - *fiber.App: The Fiber app used for routing and middleware.
func (fs *FiberServer) GetRouter() *fiber.App {
	return fs.app
}

// GracefulShutdown shuts down the server gracefully on receiving an interrupt signal.
//
// This method ensures ongoing requests are completed before shutting down.
func (fs *FiberServer) GracefulShutdown() {
	// Create a channel to listen for OS interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Shutdown the server gracefully
	if err := fs.app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown successfully")
}

// Example usage:
//
//	func main() {
//	    config := gopherfiber.ServerConfig{
//	        Port:        8080,
//	        UseTLS:      true,
//	        TLSCertFile: "./certs/server.crt",
//	        TLSKeyFile:  "./certs/server.key",
//	        UseCORS:     true,
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
