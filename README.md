Here's the detailed documentation in markdown format for your repository, grouped by packages:

---

# MyGopher

**MyGopher** is a collection of utility packages written in Go, designed to simplify common tasks such as web server setup, middleware integration, database connections (MongoDB, PostgreSQL), and token management (JWT and Paseto). Each package is modular, allowing you to integrate them into your project as needed.

## Table of Contents

- [GopherGin](#gophergin)
- [GopherMiddleware](#gophermiddleware)
- [GopherMongo](#gophermongo)
- [GopherPostgres](#gopherpostgres)
- [GopherToken](#gophertoken)

---

## GopherGin

**Package:** `gophergin`

### About

The `gophergin` package provides a set of tools for setting up a Gin web server with customizable configurations. It supports both CORS-enabled and non-CORS setups, static file serving, template loading, TLS (HTTPS) support, and graceful shutdowns. This package helps developers create production-ready HTTP servers with ease.

## Installation

To install and use the `gophergin` package, follow these steps:

1. First, make sure you have [Go](https://golang.org/dl/) installed and properly set up.
2. Run the following command to install the `gophergin` package and its dependencies:

```bash
go get github.com/lordofthemind/mygopher/gophergin
```

This command fetches the package from the repository and installs it in your Go workspace.

3. After installing, you can import and use the package in your Go project like this:

```go
import "github.com/lordofthemind/mygopher/gophergin"
```

### Types

#### `ServerConfig`
`ServerConfig` holds the configuration for setting up the Gin server.

**Fields:**
- `Port`: The port on which the server will listen (e.g., `8080`).
- `StaticPath`: Path to static files to be served (e.g., `./static`).
- `TemplatePath`: Path to HTML templates (e.g., `./templates`).
- `UseTLS`: A boolean flag to indicate whether TLS/HTTPS should be used.
- `TLSCertFile`: Path to the TLS certificate file.
- `TLSKeyFile`: Path to the TLS private key file.
- `UseCORS`: A boolean flag to enable CORS.
- `CORSConfig`: Configuration for the CORS middleware if `UseCORS` is true.

#### `ServerSetup`
`ServerSetup` is an interface for setting up a Gin server.

### Functions

#### `SetUpServer(config ServerConfig) (*gin.Engine, error)`

**CorsServerSetup**

This method sets up a Gin server with optional CORS support, static file serving, and template loading.

**Parameters:**
- `config`: An instance of `ServerConfig` with the server's configuration.

**Returns:**
- `*gin.Engine`: The Gin engine (router) instance, which can be used to define routes and start the server.
- `error`: An error if the server setup fails.

**Example usage:**

```go
serverConfig := ServerConfig{
    StaticPath:   "./static",
    TemplatePath: "./templates",
    UseCORS:      true,
    CORSConfig: cors.Config{
        AllowOrigins: []string{"https://example.com"},
        AllowMethods: []string{"GET", "POST"},
    },
}

router, err := (&CorsServerSetup{}).SetUpServer(serverConfig)
if err != nil {
    log.Fatalf("Failed to set up server: %v", err)
}
```

#### `StartGinServer(router *gin.Engine, config ServerConfig) error`

Starts the Gin server using the provided router configuration. It supports both HTTP and HTTPS (TLS) setups.

**Parameters:**
- `router`: The `*gin.Engine` instance to start.
- `config`: The `ServerConfig` containing server options like port and TLS files.

**Returns:**
- `error`: An error if the server fails to start.

**Example usage:**

```go
err := StartGinServer(router, serverConfig)
if err != nil {
    log.Fatalf("Failed to start server: %v", err)
}
```

#### `GracefulShutdown(server *http.Server)`

Handles the graceful shutdown of the server when an interrupt signal (Ctrl+C) is received. It ensures all in-flight requests are completed before shutting down within a given timeout (5 seconds).

**Parameters:**
- `server`: The `*http.Server` instance representing the running Gin server.

**Example usage:**

```go
go GracefulShutdown(&http.Server{Addr: ":8080"})
```

#### `LoadTLSCertificate(certFile, keyFile string) (tls.Certificate, error)`

Loads the TLS certificate and private key from the specified files to enable HTTPS for the server.

**Parameters:**
- `certFile`: Path to the TLS certificate file.
- `keyFile`: Path to the TLS key file.

**Returns:**
- `tls.Certificate`: The loaded TLS certificate.
- `error`: An error if the certificate fails to load.

**Example usage:**

```go
cert, err := LoadTLSCertificate("/path/to/cert.crt", "/path/to/key.key")
if err != nil {
    log.Fatalf("Failed to load TLS certificate: %v", err)
}
```

### Example Usage

```go
package main

import (
    "log"
    "github.com/gin-contrib/cors"
    "github.com/lordofthemind/mygopher/gophergin"
)

func main() {
    // Define server configuration
    serverConfig := gophergin.ServerConfig{
        Port:         8080,
        StaticPath:   "./static",
        TemplatePath: "./templates",
        UseTLS:       false,
        UseCORS:      true,
        CORSConfig: cors.Config{
            AllowOrigins: []string{"https://example.com"},
            AllowMethods: []string{"GET", "POST"},
        },
    }

    // Set up the server with CORS support
    router, err := (&gophergin.CorsServerSetup{}).SetUpServer(serverConfig)
    if err != nil {
        log.Fatalf("Failed to set up server: %v", err)
    }

    // Start the server
    if err := gophergin.StartGinServer(router, serverConfig); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }

    // Gracefully shutdown on interrupt
    gophergin.GracefulShutdown(&http.Server{Addr: ":8080"})
}
```

### Notes
- **TLS Support**: To enable HTTPS, set `UseTLS` to `true` and provide valid paths for `TLSCertFile` and `TLSKeyFile`.
- **Graceful Shutdown**: The `GracefulShutdown` function ensures that the server is terminated gracefully without abruptly closing active connections.
- **CORS Configuration**: The CORS settings can be customized via `CORSConfig` in `ServerConfig`.


---

## GopherMiddleware

**Package:** `gophermiddleware`

### About

The `gophermiddleware` package provides middleware components that can be plugged into your web server (Gin, Echo, etc.). It’s designed to handle various HTTP request and response modifications or checks.

### Example Usage

To be filled in once additional middleware features are added.

---

## GopherMongo

**Package:** `gophermongo`

### About

The `gophermongo` package simplifies the process of connecting to MongoDB, retrieving databases and collections. It supports connection retries and uses context for managing connection timeouts.

### Functions

#### `ConnectToMongoDB(ctx, dsn, timeout, maxRetries)`
Establishes a connection to MongoDB with retries and a context timeout.

**Parameters:**
- `ctx`: Context for connection management.
- `dsn`: MongoDB connection string (Data Source Name).
- `timeout`: Duration for the connection attempt.
- `maxRetries`: Maximum number of retries before giving up.

**Returns:**
- `*mongo.Client`: The connected MongoDB client instance on success.
- `error`: An error if the connection fails.

#### `GetDatabase(client, dbName)`
Retrieves a specified MongoDB database instance from the client.

**Parameters:**
- `client`: MongoDB client instance.
- `dbName`: Name of the database to retrieve.

**Returns:**
- `*mongo.Database`: The MongoDB database instance.

#### `GetCollection(db, collectionName)`
Retrieves a specified MongoDB collection from the database.

**Parameters:**
- `db`: MongoDB database instance.
- `collectionName`: Name of the collection to retrieve.

**Returns:**
- `*mongo.Collection`: The MongoDB collection instance.

### Example Usage

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/lordofthemind/mygopher/gophermongo"
)

func main() {
	ctx := context.Background()
	client, err := gophermongo.ConnectToMongoDB(ctx, "mongodb://localhost:27017", 10*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	database := gophermongo.GetDatabase(client, "myDatabase")
	collection := gophermongo.GetCollection(database, "myCollection")
	// Use collection for CRUD operations...
}
```

---

## GopherPostgres

**Package:** `gopherpostgres`

### About

The `gopherpostgres` package helps in connecting to PostgreSQL using either the standard SQL driver or GORM (an ORM library). It includes retry logic and context-based timeouts for connection management.

### Functions

#### `ConnectPostgresDB(ctx, dsn, timeout, maxRetries)`
Connects to a PostgreSQL database using the standard SQL package with retries.

**Parameters:**
- `ctx`: Context for connection management.
- `dsn`: PostgreSQL connection string (Data Source Name).
- `timeout`: Duration for the connection attempt.
- `maxRetries`: Maximum number of retries before giving up.

**Returns:**
- `*sql.DB`: The connected PostgreSQL database instance.
- `error`: An error if the connection fails.

#### `ConnectToPostgresGORM(ctx, dsn, timeout, maxRetries)`
Connects to a PostgreSQL database using GORM with retries.

**Parameters:**
- `ctx`: Context for connection management.
- `dsn`: PostgreSQL connection string (Data Source Name).
- `timeout`: Duration for the connection attempt.
- `maxRetries`: Maximum number of retries before giving up.

**Returns:**
- `*gorm.DB`: The connected GORM PostgreSQL database instance.
- `error`: An error if the connection fails.

### Example Usage

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/lordofthemind/mygopher/gopherpostgres"
)

func main() {
	ctx := context.Background()
	db, err := gopherpostgres.ConnectPostgresDB(ctx, "postgres://user:password@localhost:5432/mydb", 10*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Use db for SQL operations...
}
```

---

## GopherToken

**Package:** `gophertoken`

### About

The `gophertoken` package provides an interface and implementations for JWT and Paseto token creation and validation. It includes a payload structure and methods for generating tokens with expiration and validating them.

### Functions

#### `NewTokenManager(tokenType, secretKey)`
Creates a new token manager (JWT or Paseto) depending on the provided type.

**Parameters:**
- `tokenType`: Either "jwt" or "paseto".
- `secretKey`: Secret key for signing tokens.

**Returns:**
- `TokenManager`: An instance of the token manager.

#### `GenerateToken(username, duration)`
Generates a token for a given user and duration.

**Parameters:**
- `username`: The username for the token.
- `duration`: The token's validity duration.

**Returns:**
- `string`: The generated token.
- `error`: An error if token generation fails.

#### `ValidateToken(token)`
Validates the given token and returns the payload.

**Parameters:**
- `token`: The token to validate.

**Returns:**
- `*Payload`: The decoded token payload.
- `error`: An error if token validation fails.

### Example Usage (JWT)

```go
package main

import (
	"log"
	"time"

	"github.com/lordofthemind/mygopher/gophertoken"
)

func main() {
	manager, err := gophertoken.NewTokenManager("jwt", "your-secret-key")
	if err != nil {
		log.Fatalf("Failed to create JWT manager: %v", err)
	}

	token, err := manager.GenerateToken("user123", time.Hour)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	log.Println("Generated Token:", token)

	payload, err := manager.ValidateToken(token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	log.Printf("Token Payload: %+v\n", payload)
}
```

---

### License

This repository is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.

---

### Contributing

Feel free to contribute by submitting issues or pull requests.

---

This documentation provides a high-level overview and detailed usage examples for each package in the **MyGopher** repository. You can expand and adapt it as the project evolves!