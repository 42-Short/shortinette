package db

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

//TODO: support transactions

// Data Access Object for interacting with the DB
type DAO[T any] struct {
	DB *DB
	md metadata
}

type metadata struct {
	dbTags      []string
	primaryKeys []string
	tableName   string
}

var metadataCache sync.Map

func NewDAO[T any](db *DB) *DAO[T] {
	md := getOrExtractMetadata[T]()
	return &DAO[T]{
		DB: db,
		md: md,
	}
}

// Adds a new record to the table
func (dao *DAO[T]) Insert(ctx context.Context, data *T) error {
	query := buildInsertQuery(dao.md.tableName, dao.md.dbTags)
	_, err := dao.DB.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		return fmt.Errorf("failed to insert data in table %s: %v", dao.md.tableName, err)
	}
	return nil
}

// Modifies an existing record in the table using the DAO's primary keys.
func (dao *DAO[T]) Update(ctx context.Context, data *T) error {
	query := buildUpdateQuery(dao.md.tableName, dao.md.dbTags, dao.md.primaryKeys)
	_, err := dao.DB.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		return fmt.Errorf("failed to update data in table %s: %v", dao.md.tableName, err)
	}
	return nil
}

// Retrieves all records from the table corresponding to the DAO's type.
func (dao *DAO[T]) GetAll(ctx context.Context) ([]T, error) {
	query := buildSelectQuery(dao.md.tableName, []string{})
	var retrievedData []T
	err := dao.DB.Conn.SelectContext(ctx, &retrievedData, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.md.tableName, err)
	}
	return retrievedData, nil
}

// Retrieves a single record by the primary keys from the table.
func (dao *DAO[T]) Get(ctx context.Context, args ...any) (*T, error) {
	query := buildSelectQuery(dao.md.tableName, dao.md.primaryKeys)
	var retrievedData T
	err := dao.DB.Conn.GetContext(ctx, &retrievedData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.md.tableName, err)
	}
	return &retrievedData, err
}

// Retrieves records from the table that match the given filters.
func (dao *DAO[T]) GetFiltered(ctx context.Context, filters map[string]any) ([]T, error) {
	fields, args := extractFieldsAndArgs(filters)
	query := buildSelectQuery(dao.md.tableName, fields)
	var retrievedData []T
	err := dao.DB.Conn.SelectContext(ctx, &retrievedData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.md.tableName, err)
	}
	return retrievedData, nil
}

// Removes a record from the table using the DAO's primary keys.
func (dao *DAO[T]) Delete(ctx context.Context, args ...any) error {
	query := buildDeleteQuery(dao.md.tableName, dao.md.primaryKeys)
	_, err := dao.DB.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query on table %s: %v", dao.md.tableName, err)
	}
	return nil
}

// Example query:
// INSERT INTO participant (intra_login, github_login)
// VALUES (:intra_login, :github_login);
func buildInsertQuery(tableName string, dbTags []string) string {
	columns := strings.Join(dbTags, ", ")
	placeholders := ":" + strings.Join(dbTags, ", :")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, columns, placeholders)
}

// Example query:
//
//	UPDATE participant
//	SET intra_login = :intra_login, github_login = :github_login
//	WHERE intra_login = ? AND github_login = ?;
func buildUpdateQuery(tableName string, dbTags, primaryKeys []string) string {
	setClauses := strings.Join(buildClauses(dbTags), ", ")
	whereClauses := strings.Join(buildClauses(primaryKeys), " AND ")
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setClauses, whereClauses)
}

// Example query:
// SELECT * FROM participant
// WHERE intra_login = ? AND github_login = ?
func buildSelectQuery(tableName string, fields []string) string {
	if len(fields) == 0 {
		return fmt.Sprintf("SELECT * FROM %s", tableName)
	}
	conditions := buildConditions(fields)
	return fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, strings.Join(conditions, " AND "))
}

// Example query:
// DELETE FROM participant
// WHERE intra_login = ? AND github_login = ?
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

func getOrExtractMetadata[T any]() metadata {
	var dummy T
	typ := reflect.TypeOf(dummy)

	cached, ok := metadataCache.Load(typ)
	if ok {
		return cached.(metadata)
	}

	dbTags := extractTags(&dummy, "db")
	if len(dbTags) == 0 {
		panic("getOrExtractMetadata: Expected database tags (db) but got 0")
	}
	primaryKeys := extractTags(&dummy, "primaryKey")
	if len(primaryKeys) == 0 {
		panic("getOrExtractMetadata: Expected primaryKey tags (primaryKey) but got 0")
	}
	tableName := deriveSchemaNameFromStruct(dummy)

	md := metadata{
		dbTags:      dbTags,
		primaryKeys: primaryKeys,
		tableName:   tableName,
	}
	metadataCache.Store(typ, md)
	return md
}

func extractTags[T any](data *T, key string) []string {
	t := reflect.TypeOf(*data)

	tags := make([]string, 0, t.NumField())
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
