package database

import (
	"GoWeb/app"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"reflect"
)

// Migrate given a dummy object of any type, it will create a table with the same name as the type and create columns with the same name as the fields of the object
func Migrate(app *app.App, anyStruct interface{}) error {
	valueOfStruct := reflect.ValueOf(anyStruct)
	typeOfStruct := valueOfStruct.Type()

	tableName := typeOfStruct.Name()
	err := createTable(app, tableName)
	if err != nil {
		return err
	}

	for i := 0; i < valueOfStruct.NumField(); i++ {
		fieldType := typeOfStruct.Field(i)
		fieldName := fieldType.Name
		if fieldName != "Id" && fieldName != "id" {
			err := createColumn(app, tableName, fieldName, fieldType.Type.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// createTable creates a table with the given name if it doesn't exist, it is assumed that id will be the primary key
func createTable(app *app.App, tableName string) error {
	sanitizedTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (\"Id\" serial primary key)", tableName)

	_, err := app.Db.Query(sanitizedTableQuery)
	if err != nil {
		log.Println("Error creating table: " + tableName)
		return err
	}

	log.Println("Table created successfully (or already exists): " + tableName)
	return nil
}

// createColumn creates a column with the given name and type if it doesn't exist
func createColumn(app *app.App, tableName, columnName, columnType string) error {
	postgresType, err := getPostgresType(columnType)
	if err != nil {
		log.Println("Error creating column: " + columnName + " in table: " + tableName + " with type: " + postgresType)
		return err
	}

	sanitizedTableName := pq.QuoteIdentifier(tableName)
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS \"%s\" %s", sanitizedTableName, columnName, postgresType)

	_, err = app.Db.Query(query)
	if err != nil {
		log.Println("Error creating column: " + columnName + " in table: " + tableName + " with type: " + postgresType)
		return err
	}

	log.Println("Column created successfully (or already exists):", columnName)

	return nil
}

// Given a type in Go, return the corresponding type in Postgres
func getPostgresType(goType string) (string, error) {
	switch goType {
	case "int", "int32", "uint", "uint32":
		return "integer", nil
	case "int64", "uint64":
		return "bigint", nil
	case "int16", "int8", "uint16", "uint8", "byte":
		return "smallint", nil
	case "string":
		return "text", nil
	case "float64":
		return "double precision", nil
	case "bool":
		return "boolean", nil
	case "time.Time":
		return "timestamp", nil
	case "[]byte":
		return "bytea", nil
	}

	return "", errors.New("Unknown type: " + goType)
}
