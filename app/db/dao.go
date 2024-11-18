package db

import (
	"fmt"
	"reflect"
	"strings"
)

type DAO[T any] struct {
	DB          *DB
	tableName   string
	dbTags      []string
	primaryKeys []string
}

func NewDAO[T any](db *DB) *DAO[T] {
	var dummy T
	tags := extractTags(&dummy, "db")
	if len(tags) == 0 {
		panic("NewDAO: Expected database tags (db) but got 0")
	}
	primaryKeys := extractTags(&dummy, "primaryKey")
	if len(primaryKeys) == 0 {
		panic("NewDAO: Expected primaryKey tags (primaryKey) but got 0")
	}
	tableName := deriveSchemaNameFromStruct(dummy)
	return &DAO[T]{
		DB:          db,
		tableName:   tableName,
		dbTags:      tags,
		primaryKeys: primaryKeys,
	}
}

// Adds a new record to the table.
func (dao *DAO[T]) Insert(data *T) error {
	query := buildInsertQuery(dao.tableName, dao.dbTags)
	_, err := dao.DB.namedExecWithTimeout(query, data)
	if err != nil {
		return fmt.Errorf("failed to insert data in table %s: %v", dao.tableName, err)
	}
	return nil
}

// Modifies an existing record in the table using the DAO's primary keys.
func (dao *DAO[T]) Update(data *T) error {
	query := buildUpdateQuery(dao.tableName, dao.dbTags, dao.primaryKeys)
	_, err := dao.DB.namedExecWithTimeout(query, data)
	if err != nil {
		return fmt.Errorf("failed to update data in table %s: %v", dao.tableName, err)
	}
	return nil
}

// Retrieves all records from the table corresponding to the DAO's type.
func (dao *DAO[T]) GetAll() ([]T, error) {
	query := buildSelectQuery(dao.tableName, []string{})
	var retrievedData []T
	err := dao.DB.selectWithTimeout(&retrievedData, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}
	return retrievedData, nil
}

// Retrieves a single record by the primary keys from the table.
func (dao *DAO[T]) Get(args ...any) (*T, error) {
	query := buildSelectQuery(dao.tableName, dao.primaryKeys)
	var retrievedData T
	err := dao.DB.getWithTimeout(&retrievedData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}
	return &retrievedData, err
}

// Retrieves records from the table that match the given filters.
func (dao *DAO[T]) GetFiltered(filters map[string]any) ([]T, error) {
	fields, args := extractFieldsAndArgs(filters)
	query := buildSelectQuery(dao.tableName, fields)
	var retrievedData []T
	err := dao.DB.selectWithTimeout(&retrievedData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}
	return retrievedData, nil
}

// Removes a record from the table using the DAO's primary keys.
func (dao *DAO[T]) Delete(args ...any) error {
	query := buildDeleteQuery(dao.tableName, dao.primaryKeys)
	_, err := dao.DB.execWithTimeout(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query on table %s: %v", dao.tableName, err)
	}
	return nil
}

func buildInsertQuery(tableName string, dbTags []string) string {
	columns := strings.Join(dbTags, ", ")
	placeholders := ":" + strings.Join(dbTags, ", :")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, columns, placeholders)
}

func buildUpdateQuery(tableName string, dbTags, primaryKeys []string) string {
	setClauses := strings.Join(buildClauses(dbTags), ", ")
	whereClauses := strings.Join(buildClauses(primaryKeys), " AND ")
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setClauses, whereClauses)
}

func buildSelectQuery(tableName string, fields []string) string {
	if len(fields) == 0 {
		return fmt.Sprintf("SELECT * FROM %s", tableName)
	}
	conditions := buildConditions(fields)
	return fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, strings.Join(conditions, " AND "))
}

func buildDeleteQuery(tableName string, primaryKeys []string) string {
	conditions := buildConditions(primaryKeys)
	return fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, strings.Join(conditions, " AND "))
}

func buildClauses(fields []string) []string {
	clauses := make([]string, len(fields))
	for i, field := range fields {
		clauses[i] = fmt.Sprintf("%s = :%s", field, field)
	}
	return clauses
}

func buildConditions(fields []string) []string {
	conditions := make([]string, len(fields))
	for i, field := range fields {
		conditions[i] = fmt.Sprintf("%s = ?", field)
	}
	return conditions
}

func extractTags[T any](data *T, key string) []string {
	t := reflect.TypeOf(*data)
	tags := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(key)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

func extractFieldsAndArgs(filters map[string]any) ([]string, []any) {
	fields := make([]string, 0, len(filters))
	args := make([]any, 0, len(filters))
	for field, value := range filters {
		fields = append(fields, field)
		args = append(args, value)
	}
	return fields, args
}
