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

The `gophergin` package provides a simple function to set up a Gin web server. It’s designed to ease the process of setting up routes and middleware in a clean and maintainable way.

### Functions

#### `SetUpGinServer()`
Sets up and returns a configured Gin server with standard settings.

**Returns:**
- `*gin.Engine`: A new Gin web server instance.

### Example Usage

```go
package main

import (
	"github.com/lordofthemind/mygopher/gophergin"
)

func main() {
	server := gophergin.SetUpGinServer()
	server.Run(":8080")
}
```

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