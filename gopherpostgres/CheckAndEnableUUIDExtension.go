package gopherpostgres

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// CheckAndEnableUUIDExtension checks if the 'uuid-ossp' extension is enabled in PostgreSQL and enables it if it is not.
//
// The 'uuid-ossp' extension is required for generating UUIDs using the `uuid_generate_v4()` function in PostgreSQL.
// This function queries the PostgreSQL system catalog to check whether the extension is already enabled. If the extension
// is not found, it attempts to create and enable the extension.
//
// Params:
//
//	db - The GORM database connection (*gorm.DB) used to interact with the PostgreSQL database.
//
// Returns:
//
//	error - Returns an error if the check or enabling the extension fails.
//
// Example usage:
//
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//	    log.Fatalf("Failed to connect to PostgreSQL: %v", err)
//	}
//
//	err = CheckAndEnableUUIDExtension(db)
//	if err != nil {
//	    log.Fatalf("Error enabling UUID extension: %v", err)
//	}
//
// This function is useful when working with GORM and PostgreSQL to ensure the 'uuid-ossp' extension is available
// for UUID generation in your database schema.
func CheckAndEnableUUIDExtension(db *gorm.DB) error {
	// Get the underlying sql.DB connection from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection from GORM: %v", err)
	}

	// Query PostgreSQL to check if the 'uuid-ossp' extension is already enabled
	row := sqlDB.QueryRow("SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp'")
	var exists int
	err = row.Scan(&exists)

	// Handle the case where the extension is not found
	if err == sql.ErrNoRows {
		// Attempt to enable the 'uuid-ossp' extension
		_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		if err != nil {
			return fmt.Errorf("failed to create uuid-ossp extension: %v", err)
		}
		log.Println("uuid-ossp extension enabled successfully")
	} else if err != nil {
		// Return an error if the query failed for any other reason
		return fmt.Errorf("failed to check for uuid-ossp extension: %v", err)
	} else {
		// Log success if the extension is already enabled
		log.Println("uuid-ossp extension is already enabled")
	}

	return nil
}
