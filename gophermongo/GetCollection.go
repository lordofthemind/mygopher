package gophermongo

import "go.mongodb.org/mongo-driver/mongo"

// GetCollection retrieves the specified MongoDB collection from the database.
//
// This function allows you to get a reference to a specific collection within a MongoDB
// database. Once you have the collection, you can perform CRUD operations on it.
//
// Params:
//
//	db - The MongoDB database instance.
//	collectionName - The name of the collection to retrieve.
//
// Returns:
//
//	*mongo.Collection - The MongoDB collection instance.
//
// Example usage:
//
//	database := GetDatabase(client, "myDatabase")
//	collection := GetCollection(database, "myCollection")
//
// After retrieving the collection, you can perform operations like insert, find, update, or delete.
func GetCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}
