Here's a detailed documentation guide for each package in your repository. This guide includes an overview, purpose, installation instructions, and usage examples for each package, making it easy for other developers to understand and use your code.

---

## **gophergin** Documentation

### **Overview**

`gophergin` provides utilities for setting up a Gin server with optional CORS and TLS support. It abstracts the configuration of CORS policies and allows the creation of secure servers with minimal effort.

### **Installation**

To use this package, you can install it using the `go get` command:

```bash
go get github.com/yourusername/repositoryname/gophergin
```

### **Functions**

#### **SetUpServerWithOptionalCORS()**

This function sets up a Gin server with or without CORS based on a configuration. It returns a `*gin.Engine` instance.

- **Example**:
    ```go
    router, err := gophergin.SetUpServerWithOptionalCORS()
    if err != nil {
        log.Fatal("Failed to set up server:", err)
    }
    ```

#### **StartGinServer(router *gin.Engine)**

Starts the Gin server on a specified port, with TLS if enabled.

- **Example**:
    ```go
    err := gophergin.StartGinServer(router)
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
    ```

#### **GracefulShutdown(server *http.Server)**

Handles the graceful shutdown of the Gin server.

- **Example**:
    ```go
    gophergin.GracefulShutdown(server)
    ```

### **Usage Example**

```go
package main

import (
    "log"
    "github.com/yourusername/repositoryname/gophergin"
)

func main() {
    router, err := gophergin.SetUpServerWithOptionalCORS()
    if err != nil {
        log.Fatal("Error setting up server:", err)
    }

    // Start the server
    err = gophergin.StartGinServer(router)
    if err != nil {
        log.Fatal("Error starting server:", err)
    }
}
```

---

## **gophermongo** Documentation

### **Overview**

`gophermongo` provides utilities for connecting to MongoDB and retrieving collections and databases. It simplifies MongoDB setup by providing connection helpers.

### **Installation**

```bash
go get github.com/yourusername/repositoryname/gophermongo
```

### **Functions**

#### **ConnectMongoDB(uri string) (*mongo.Client, error)**

Connects to a MongoDB instance and returns a Mongo client.

- **Example**:
    ```go
    client, err := gophermongo.ConnectMongoDB("mongodb://localhost:27017")
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    ```

#### **GetDatabase(client *mongo.Client, dbName string) (*mongo.Database, error)**

Retrieves a MongoDB database from the client.

- **Example**:
    ```go
    db, err := gophermongo.GetDatabase(client, "mydb")
    if err != nil {
        log.Fatal("Failed to get database:", err)
    }
    ```

#### **GetCollection(db *mongo.Database, collectionName string) (*mongo.Collection, error)**

Retrieves a specific collection from a MongoDB database.

- **Example**:
    ```go
    collection, err := gophermongo.GetCollection(db, "users")
    if err != nil {
        log.Fatal("Failed to get collection:", err)
    }
    ```

### **Usage Example**

```go
package main

import (
    "log"
    "github.com/yourusername/repositoryname/gophermongo"
)

func main() {
    // Connect to MongoDB
    client, err := gophermongo.ConnectMongoDB("mongodb://localhost:27017")
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }

    // Get Database
    db, err := gophermongo.GetDatabase(client, "mydb")
    if err != nil {
        log.Fatal("Failed to get database:", err)
    }

    // Get Collection
    collection, err := gophermongo.GetCollection(db, "users")
    if err != nil {
        log.Fatal("Failed to get collection:", err)
    }

    log.Println("Successfully connected to collection:", collection.Name())
}
```

---

## **gopherpostgres** Documentation

### **Overview**

`gopherpostgres` provides utilities to connect to PostgreSQL databases using both `database/sql` and GORM. It simplifies the connection setup to PostgreSQL.

### **Installation**

```bash
go get github.com/yourusername/repositoryname/gopherpostgres
```

### **Functions**

#### **ConnectPostgresDB(dsn string) (*sql.DB, error)**

Connects to PostgreSQL using the provided DSN (Data Source Name) and returns a `*sql.DB` instance.

- **Example**:
    ```go
    db, err := gopherpostgres.ConnectPostgresDB("host=localhost user=postgres password=mysecretpassword dbname=mydb sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    ```

#### **ConnectPostgresGORM(dsn string) (*gorm.DB, error)**

Connects to PostgreSQL using GORM and returns a `*gorm.DB` instance.

- **Example**:
    ```go
    gormDB, err := gopherpostgres.ConnectPostgresGORM("host=localhost user=postgres password=mysecretpassword dbname=mydb sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL via GORM:", err)
    }
    ```

### **Usage Example**

```go
package main

import (
    "log"
    "github.com/yourusername/repositoryname/gopherpostgres"
)

func main() {
    // Connect to PostgreSQL using database/sql
    db, err := gopherpostgres.ConnectPostgresDB("host=localhost user=postgres password=mysecretpassword dbname=mydb sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    log.Println("Connected to PostgreSQL")

    // Connect to PostgreSQL using GORM
    gormDB, err := gopherpostgres.ConnectPostgresGORM("host=localhost user=postgres password=mysecretpassword dbname=mydb sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL via GORM:", err)
    }
    log.Println("Connected to PostgreSQL using GORM:", gormDB.Name())
}
```

---

## **gophertoken** Documentation

### **Overview**

`gophertoken` provides an easy-to-use interface for generating and validating tokens using both JWT and Paseto encryption. It allows you to choose between token types, with built-in support for symmetric keys.

### **Installation**

```bash
go get github.com/yourusername/repositoryname/gophertoken
```

### **Interfaces**

#### **TokenManager Interface**

```go
type TokenManager interface {
    GenerateToken(username string, duration time.Duration) (string, error)
    ValidateToken(token string) (*Payload, error)
}
```

#### **NewTokenManager(tokenType, secretKey string) (TokenManager, error)**

This function returns a `TokenManager` that can handle either JWT or Paseto tokens based on the `tokenType` parameter.

- **Parameters**:
    - `tokenType`: Choose either `"jwt"` or `"paseto"`.
    - `secretKey`: Symmetric key for signing tokens.

- **Example**:
    ```go
    tokenManager, err := gophertoken.NewTokenManager("jwt", "your-secret-key")
    if err != nil {
        log.Fatal("Failed to create token manager:", err)
    }
    ```

### **Functions**

#### **GenerateToken(username string, duration time.Duration) (string, error)**

Generates a token for a user with a specific validity duration.

- **Example**:
    ```go
    token, err := tokenManager.GenerateToken("john_doe", time.Hour)
    if err != nil {
        log.Fatal("Failed to generate token:", err)
    }
    ```

#### **ValidateToken(token string) (*Payload, error)**

Validates a token and returns its payload (username, expiration, etc.).

- **Example**:
    ```go
    payload, err := tokenManager.ValidateToken(token)
    if err != nil {
        log.Fatal("Invalid token:", err)
    }
    ```

### **Usage Example**

```go
package main

import (
    "log"
    "github.com/yourusername/repositoryname/gophertoken"
    "time"
)

func main() {
    // Create a token manager (JWT)
    tokenManager, err := gophertoken.NewTokenManager("jwt", "your-secret-key")
    if err != nil {
        log.Fatal("Error creating token manager:", err)
    }

    // Generate a token
    token, err := tokenManager.GenerateToken("john_doe", time.Hour)
    if err != nil {
        log.Fatal("Error generating token:", err)
    }
    log.Println("Generated Token:", token)

    // Validate the token
    payload, err := tokenManager.ValidateToken(token)
    if err != nil {
        log.Fatal("Error validating token:", err)
    }
    log.Println("Token valid for user:", payload.Username)
}
```

---

## **General Notes**

- **Modular Design**: Each package in the repository is designed to be self-contained and reusable. This makes it easy to include them in various projects.
- **Security**: Be

 sure to use strong, unique secret keys in production environments for the `gophertoken` package.
- **Usage Examples**: Each package includes a detailed usage example, which can be quickly referenced in your own applications.

By following this documentation, other developers can easily integrate and use the functionality provided by each package in your repository.