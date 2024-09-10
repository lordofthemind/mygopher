# mygopher - A Go Utilities Package

`mygopher` is a modular Go package that provides utility functions for connecting to MongoDB and PostgreSQL databases, managing Gin servers, and setting up logging. This package is designed to simplify backend development by offering reusable components.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Set Up Logger](#set-up-logger)
  - [Connect to MongoDB](#connect-to-mongodb)
  - [Connect to PostgreSQL](#connect-to-postgresql)
  - [Connect to PostgreSQL using GORM](#connect-to-postgresql-using-gorm)
- [Examples](#examples)
- [License](#license)
- [Contributing](#contributing)
- [Contact](#contact)

## Installation

Install the package using Go modules:

```bash
go get github.com/lordofthemind/mygopher@v0.1.0
```

Replace `v0.1.0` with the desired version.

## Usage

Below are the primary functionalities provided by the MyGopher package:

### Set Up Logger

Set up logging to both a file and stdout with automatic log file management.

```go
package main

import (
    "log"
    "github.com/lordofthemind/mygopher"
)

func main() {
    logFile, err := mygopher.SetUpLoggerFile("app.log")
    if err != nil {
        log.Fatalf("Failed to set up logger: %v", err)
    }
    defer logFile.Close()
}
```

### Connect to MongoDB

Establish a connection to MongoDB with retry logic and a context-based timeout.

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/lordofthemind/mygopher"
)

func main() {
    ctx := context.Background()
    dsn := os.Getenv("MONGO_DOCKER_CONNECTION_URL")
    timeout := 30 * time.Second
    maxRetries := 5
    dbName := "polyglot" // Replace with your actual database name

    client, db, err := mygopher.ConnectToMongoDB(ctx, dsn, timeout, maxRetries, dbName)
    if err != nil {
        log.Fatalf("Error connecting to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)
}
```

### Connect to PostgreSQL

Connect to a PostgreSQL database using the `database/sql` package with context-based timeout and retry logic.

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/lordofthemind/mygopher"
)

func main() {
    ctx := context.Background()
    dsn := "postgres://username:password@localhost:5432/dbname?sslmode=disable"
    timeout := 30 * time.Second
    maxRetries := 3

    db, err := mygopher.ConnectPostgresDB(ctx, dsn, timeout, maxRetries)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer db.Close()
}
```

### Connect to PostgreSQL using GORM

Connect to a PostgreSQL database using GORM with retry logic and context-based timeout.

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/lordofthemind/mygopher"
)

func main() {
    ctx := context.Background()
    dsn := "postgres://username:password@localhost:5432/dbname?sslmode=disable"
    timeout := 30 * time.Second
    maxRetries := 3

    db, err := mygopher.ConnectToPostgreSQLGormDB(ctx, dsn, timeout, maxRetries)
    if err != nil {
        log.Fatalf("Error connecting to the database with GORM: %v", err)
    }
}
```

## Examples

Check the examples in the `/examples` directory for more detailed use cases.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For questions or suggestions, please contact [MANISH KUMAR](mailto:manish.keshav98@gmail.com).