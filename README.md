# mygopher

This Go package provides a set of utilities for working with various backend components, including PostgreSQL, MongoDB, Gin server management, and logging. It is designed to be modular and easy to integrate into your existing projects.

## Project Structure

```
.
├── LICENSE
├── README.md
├── go.mod
├── go.sum
├── main.go
├── mygopherginserver
│   ├── StartGinServer.go
│   └── StopGinServer.go
├── mygopherlogger
│   └── SetupLoggerFile.go
├── mygophermongodb
│   └── ConnectMongoDB.go
├── mygopherpostgres
│   ├── ConnectPostgresDB.go
│   └── ConnectPostgresGORM.go
└── tree.txt
```

## Installation

To use this package in your project, you can install it via Go modules:

```bash
go get github.com/lordofthemind/mygopher@v0.1.0
```

Replace `v0.1.0` with the desired version.

## Usage

### Gin Server Management (`mygopherginserver`)

This package provides functions to start and stop a Gin server.

#### Start the Gin Server

```go
package main

import (
    "github.com/lordofthemind/mygopher/mygopherginserver"
)

func main() {
    err := mygopherginserver.StartGinServer()
    if err != nil {
        log.Fatalf("Failed to start Gin server: %v", err)
    }
}
```

#### Stop the Gin Server

```go
package main

import (
    "github.com/lordofthemind/mygopher/mygopherginserver"
)

func main() {
    err := mygopherginserver.StopGinServer()
    if err != nil {
        log.Fatalf("Failed to stop Gin server: %v", err)
    }
}
```

### Logger Setup (`mygopherlogger`)

Set up logging to both a file and stdout with ease.

#### Example

```go
package main

import (
    "log"
    "github.com/lordofthemind/mygopher/mygopherlogger"
)

func main() {
    logFile, err := mygopherlogger.SetupLoggerFile("app.log")
    if err != nil {
        log.Fatalf("Failed to set up logger: %v", err)
    }
    defer logFile.Close()
}
```

### MongoDB Connection (`mygophermongodb`)

Connect to a MongoDB instance with retries.

#### Example

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/lordofthemind/mygopher/mygophermongodb"
)

func main() {
    ctx := context.Background()
    dsn := "mongodb://localhost:27017"
    timeout := 30 * time.Second
    maxRetries := 3
    dbName := "mydatabase"

    client, db, err := mygophermongodb.ConnectToMongoDB(ctx, dsn, timeout, maxRetries, dbName)
    if err != nil {
        log.Fatalf("Error connecting to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)
}
```

### PostgreSQL Connection (`mygopherpostgres`)

#### Using `database/sql`

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/lordofthemind/mygopher/mygopherpostgres"
)

func main() {
    ctx := context.Background()
    dsn := "postgres://username:password@localhost:5432/dbname?sslmode=disable"
    timeout := 30 * time.Second
    maxRetries := 3

    db, err := mygopherpostgres.ConnectPostgresDB(ctx, dsn, timeout, maxRetries)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer db.Close()
}
```

#### Using GORM

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/lordofthemind/mygopher/mygopherpostgres"
)

func main() {
    ctx := context.Background()
    dsn := "postgres://username:password@localhost:5432/dbname?sslmode=disable"
    timeout := 30 * time.Second
    maxRetries := 3

    db, err := mygopherpostgres.ConnectToPostgreSQLGormDB(ctx, dsn, timeout, maxRetries)
    if err != nil {
        log.Fatalf("Error connecting to the database with GORM: %v", err)
    }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For questions or suggestions, please contact [MANISH KUMAR](mailto:manish.keshav98@gmail.com).