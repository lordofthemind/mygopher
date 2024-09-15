package gophermongo

import "go.mongodb.org/mongo-driver/mongo"

// GetDatabase retrieves the specified MongoDB database instance from the client.
//
// This function is a simple utility to get a reference to a MongoDB database using
// an existing client. It allows interaction with collections within the specified database.
//
// Params:
//
//	client - The MongoDB client instance.
//	dbName - The name of the database to retrieve.
//
// Returns:
//
//	*mongo.Database - The MongoDB database instance.
//
// Example usage:
//
//	database := GetDatabase(client, "myDatabase")
//	collection := database.Collection("myCollection")
//
// Once you have the database, you can access any collection within it.
func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}
