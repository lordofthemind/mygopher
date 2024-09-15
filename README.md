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

The `gophermongo` package provides utilities to establish a connection to a MongoDB server, retrieve databases, and access collections. It supports connection retries, context-based timeouts, and simplified database and collection retrieval.

This package is designed to handle common MongoDB operations with minimal configuration, allowing for easy integration with Go projects.

### Installation

To install and use the `gophermongo` package, run the following command:

```bash
go get github.com/lordofthemind/mygopher/gophermongo
```

After installation, import the package in your Go project:

```go
import "github.com/lordofthemind/mygopher/gophermongo"
```

### Functions

#### `ConnectToMongoDB(ctx, dsn, timeout, maxRetries)`

Establishes a connection to MongoDB with a specified context timeout and retry mechanism. The function attempts to connect up to `maxRetries` times, waiting 5 seconds between each retry.

**Parameters:**
- `ctx`: Context for managing the connection timeout.
- `dsn`: The MongoDB connection string (Data Source Name).
- `timeout`: Duration for the connection timeout (e.g., `10*time.Second`).
- `maxRetries`: Maximum number of retries before returning an error.

**Returns:**
- `*mongo.Client`: The connected MongoDB client instance on success.
- `error`: An error if the connection fails.

**Details:**
- The function will retry the connection in case of failure up to the specified number of retries (`maxRetries`).
- If a connection cannot be established within the provided `timeout`, it will return an error indicating the context has timed out.

**Example Usage:**

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

	// Now you can use the client to interact with the database
}
```

---

#### `GetDatabase(client, dbName)`

Retrieves a specified MongoDB database instance from the connected client.

**Parameters:**
- `client`: The MongoDB client instance.
- `dbName`: The name of the database to retrieve.

**Returns:**
- `*mongo.Database`: The MongoDB database instance.

**Details:**
- Use this function to retrieve a database once you have successfully connected to MongoDB.

**Example Usage:**

```go
database := gophermongo.GetDatabase(client, "myDatabase")
```

---

#### `GetCollection(db, collectionName)`

Retrieves a specified MongoDB collection from the given database.

**Parameters:**
- `db`: The MongoDB database instance.
- `collectionName`: The name of the collection to retrieve.

**Returns:**
- `*mongo.Collection`: The MongoDB collection instance.

**Details:**
- Once you have a collection, you can perform CRUD operations on it (insert, find, update, delete).

**Example Usage:**

```go
collection := gophermongo.GetCollection(database, "myCollection")
```

---

### Example Usage (Full)

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/lordofthemind/mygopher/gophermongo"
)

func main() {
	// Set up context and connection parameters
	ctx := context.Background()
	client, err := gophermongo.ConnectToMongoDB(ctx, "mongodb://localhost:27017", 10*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Get the database
	database := gophermongo.GetDatabase(client, "myDatabase")

	// Get the collection
	collection := gophermongo.GetCollection(database, "myCollection")

	// Now you can use the collection for CRUD operations
	// Example: collection.InsertOne(ctx, document)
}
```
---

## GopherPostgres

**Package:** `gopherpostgres`

### About

The `gopherpostgres` package facilitates connecting to PostgreSQL databases using either the standard SQL package or GORM (an ORM library). It includes retry logic and context-based timeouts, making it resilient and flexible for production use. This package provides functions to connect using raw SQL (`*sql.DB`) or GORM (`*gorm.DB`), making it easy to integrate with both standard and ORM-based database interactions.

### Installation

To install and use the `gopherpostgres` package, run the following command:

```bash
go get github.com/lordofthemind/mygopher/gopherpostgres
```

Then, import the package in your Go project:

```go
import "github.com/lordofthemind/mygopher/gopherpostgres"
```

### Functions

#### `ConnectPostgresDB(ctx, dsn, timeout, maxRetries)`

Connects to a PostgreSQL database using the standard SQL driver. The function includes retry logic, which attempts to reconnect in case of failure up to the specified number of retries.

**Parameters:**
- `ctx`: Context for managing connection timeout and cancellation.
- `dsn`: PostgreSQL connection string (Data Source Name).
- `timeout`: Duration for the connection attempt (e.g., `10*time.Second`).
- `maxRetries`: Maximum number of retries before returning an error.

**Returns:**
- `*sql.DB`: The connected PostgreSQL database instance on success.
- `error`: An error if the connection fails after all retries.

**Details:**
- If the connection attempt fails, the function will retry up to `maxRetries` times, with a 5-second delay between each retry.
- The function utilizes `sql.Open` and `db.PingContext` to ensure a valid connection is established.

**Example Usage:**

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

	// Perform SQL operations
	// Example: db.QueryContext(ctx, "SELECT * FROM users")
}
```

---

#### `ConnectToPostgresGORM(ctx, dsn, timeout, maxRetries)`

Connects to a PostgreSQL database using GORM. Similar to the raw SQL connection, this function includes retry logic and context-based timeouts to manage the connection lifecycle.

**Parameters:**
- `ctx`: Context for managing connection timeout and cancellation.
- `dsn`: PostgreSQL connection string (Data Source Name).
- `timeout`: Duration for the connection attempt (e.g., `10*time.Second`).
- `maxRetries`: Maximum number of retries before returning an error.

**Returns:**
- `*gorm.DB`: The connected GORM PostgreSQL database instance on success.
- `error`: An error if the connection fails after all retries.

**Details:**
- This function leverages GORM's `gorm.Open` to establish a connection using PostgreSQL. 
- If the connection fails, it retries up to `maxRetries` times with a 5-second delay between each retry.
- Ideal for projects using an ORM for database operations.

**Example Usage:**

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/lordofthemind/mygopher/gopherpostgres"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()
	db, err := gopherpostgres.ConnectToPostgresGORM(ctx, "postgres://user:password@localhost:5432/mydb", 10*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL using GORM: %v", err)
	}

	// Use GORM for ORM-based database operations
	// Example: db.Find(&users)
}
```

---

### Example Usage (Full)

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

	// Using standard SQL driver
	dbSQL, err := gopherpostgres.ConnectPostgresDB(ctx, "postgres://user:password@localhost:5432/mydb", 10*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer dbSQL.Close()

	// Using GORM
	dbGORM, err := gopherpostgres.ConnectToPostgresGORM(ctx, "postgres://user:password@localhost:5432/mydb", 10*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL using GORM: %v", err)
	}

	// Perform operations with dbSQL or dbGORM
}
```

## GopherToken

**Package:** `gophertoken`

### About

The `gophertoken` package provides a framework for token management using JWT and Paseto, supporting token generation and validation with expiration logic. It offers interfaces for both types and implements payload structures, allowing developers to manage token-based authentication.

### Payload Structure

#### `Payload`
Represents the token's claims, containing user information and expiration details.

**Fields:**
- `ID`: A unique token ID (`uuid.UUID`).
- `Username`: The username associated with the token.
- `IssuedAt`: The time when the token was issued.
- `ExpiredAt`: The expiration time of the token.

### Errors
- `ErrInvalidToken`: Indicates that the token is invalid.
- `ErrExpiredToken`: Indicates that the token has expired.

### Methods

#### `NewPayload(username, duration)`
Creates a new `Payload` for a given user and token duration.

**Parameters:**
- `username`: The username for the token.
- `duration`: Token validity period (e.g., 1 hour).

**Returns:**
- `*Payload`: The payload object containing user information and expiration data.
- `error`: If there's an issue during payload creation.

#### `Valid()`
Validates if the token has expired.

**Returns:**
- `error`: `ErrExpiredToken` if the token is expired, or `nil` if valid.

### TokenManager Interface

#### `TokenManager`
Interface for token operations like generation and validation, supporting JWT and Paseto tokens.

**Methods:**
- `GenerateToken(username string, duration time.Duration) (string, error)`: Generates a token.
- `ValidateToken(token string) (*Payload, error)`: Validates a token and returns its payload.

### Implementations

#### `NewTokenManager(tokenType, secretKey)`
Instantiates a new `TokenManager` based on the provided token type (`"jwt"` or `"paseto"`) and a secret key.

**Parameters:**
- `tokenType`: Either `"jwt"` or `"paseto"`.
- `secretKey`: The secret key used for signing the token.

**Returns:**
- `TokenManager`: A token manager that supports the specified token type.
- `error`: If the token type is unsupported or the secret key is invalid.

---

### JWT Implementation

#### `JWTMaker`
Handles JWT token creation and validation.

##### `NewJWTMaker(secretKey)`
Creates a new `JWTMaker` for JWT tokens.

**Parameters:**
- `secretKey`: Secret key for signing JWTs.

**Returns:**
- `*JWTMaker`: The JWT manager.
- `error`: If the secret key is invalid.

##### `GenerateToken(username, duration)`
Generates a JWT with a specific username and duration.

**Returns:**
- `string`: The JWT as a string.
- `error`: If there's an error during token generation.

##### `ValidateToken(tokenString)`
Validates the JWT and returns its payload.

**Returns:**
- `*Payload`: The token's payload.
- `error`: If the token is invalid or expired.

---

### Paseto Implementation

#### `PasetoMaker`
Handles Paseto token creation and validation.

##### `NewPasetoMaker(secretKey)`
Creates a new `PasetoMaker` for Paseto tokens.

**Parameters:**
- `secretKey`: Symmetric key for Paseto (must be 32 bytes).

**Returns:**
- `*PasetoMaker`: The Paseto manager.
- `error`: If the secret key is invalid.

##### `GenerateToken(username, duration)`
Generates a Paseto token with a specific username and duration.

**Returns:**
- `string`: The Paseto token.
- `error`: If token generation fails.

##### `ValidateToken(token)`
Validates the Paseto token and returns its payload.

**Returns:**
- `*Payload`: The token's payload.
- `error`: If the token is invalid or expired.

---

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

This repository is licensed under the MIT License. The full text of the license is as follows:

---

**MIT License**

Copyright (c) 2024 Manish Kumar

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

---

### Contributing

We welcome contributions to the **MyGopher** repository. To contribute:

1. **Report Issues**: If you find any bugs or have feature requests, please open an issue on the [GitHub Issues page](#).
2. **Submit Pull Requests**: For code contributions, fork the repository, make your changes, and submit a pull request. Ensure your changes adhere to the existing code style and include appropriate tests if applicable.

Please follow our [contributing guidelines](#) to ensure a smooth collaboration process.

---

### Documentation Overview

This documentation provides a comprehensive overview and detailed usage examples for each package within the **MyGopher** repository. It is intended to help developers understand the functionalities of the packages and how to integrate them into their projects effectively. The documentation is organized as follows:

- **gophergin**: Setup and configuration of the Gin server with optional CORS and TLS support.
- **gophermongo**: Utilities for connecting to MongoDB, accessing collections, and managing databases.
- **gopherpostgres**: Tools for connecting to PostgreSQL using both `database/sql` and GORM.
- **gophertoken**: Token management utilities including JWT and Paseto implementations.

As the project evolves, this documentation will be updated with new features, changes, and best practices. Feel free to contribute suggestions and improvements.