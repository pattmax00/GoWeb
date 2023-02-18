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
	// Check to see if the table already exists
	var tableExists bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_class c JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE c.relname ~ $1 AND pg_catalog.pg_table_is_visible(c.oid))", "^"+tableName+"$").Scan(&tableExists)
	if err != nil {
		log.Println("Error checking if table exists: " + tableName)
		return err
	}

	if tableExists {
		log.Println("Table already exists: " + tableName)
		return nil
	} else {
		sanitizedTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (\"Id\" serial primary key)", tableName)

		_, err := app.Db.Query(sanitizedTableQuery)
		if err != nil {
			log.Println("Error creating table: " + tableName)
			return err
		}

		log.Println("Table created successfully: " + tableName)
		return nil
	}
}

// createColumn creates a column with the given name and type if it doesn't exist
func createColumn(app *app.App, tableName, columnName, columnType string) error {
	// Check to see if the column already exists
	var columnExists bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", tableName, columnName).Scan(&columnExists)
	if err != nil {
		log.Println("Error checking if column exists: " + columnName + " in table: " + tableName)
		return err
	}

	if columnExists {
		log.Println("Column already exists: " + columnName + " in table: " + tableName)
		return nil
	} else {
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

		log.Println("Column created successfully:", columnName)

		return nil
	}
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
	case "Time":
		return "timestamp", nil
	case "[]byte":
		return "bytea", nil
	}

	return "", errors.New("Unknown type: " + goType)
}
