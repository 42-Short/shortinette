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

//Update Method
// Delete Method

func NewDAO[T any](db *DB) *DAO[T] {
	var dummy T
	tags := extractTags(&dummy, "db")
	if len(tags) == 0 {
		panic("NewBaseDao: Expected database tags (db) but got 0")
	}
	primaryKeys := extractTags(&dummy, "primaryKey")
	if len(primaryKeys) == 0 {
		panic("NewBaseDao: Expected primaryKey tags (primaryKey) but got 0")
	}
	tableName := deriveSchemaNameFromStruct(dummy)
	return &DAO[T]{
		DB:          db,
		tableName:   tableName,
		dbTags:      tags,
		primaryKeys: primaryKeys,
	}
}

func (dao *DAO[T]) Insert(data *T) error {
	columns := strings.Join(dao.dbTags, ", ")
	placeholders := ":" + strings.Join(dao.dbTags, ", :")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", dao.tableName, columns, placeholders)

	_, err := dao.DB.namedExecWithTimeout(query, data)
	if err != nil {
		return fmt.Errorf("failed to insert data into table %s: %v", dao.tableName, err)
	}
	return nil
}

func (dao *DAO[T]) GetAll() ([]T, error) {
	query := dao.buildSelectQuery([]string{})

	var retrievedData []T
	err := dao.DB.selectWithTimeout(&retrievedData, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}
	return retrievedData, nil
}

func (dao *DAO[T]) Get(args ...any) (*T, error) {
	var retrievedData T

	query := dao.buildSelectQuery(dao.primaryKeys)
	err := dao.DB.getWithTimeout(&retrievedData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}
	return &retrievedData, err
}

func (dao *DAO[T]) GetFiltered(filters map[string]any) ([]T, error) {
	fields := make([]string, 0, len(filters))
	args := make([]any, 0, len(filters))
	for field, value := range filters {
		fields = append(fields, field)
		args = append(args, value)
	}

	var retrievedData []T
	query := dao.buildSelectQuery(fields)
	err := dao.DB.selectWithTimeout(&retrievedData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter data: %v", err)
	}

	return retrievedData, nil
}

func (dao *DAO[T]) Delete(args ...any) error {
	conditions := buildConditions(dao.primaryKeys)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", dao.tableName, strings.Join(conditions, " AND "))
	_, err := dao.DB.execWithTimeout(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete from table %s: %v", dao.tableName, err)
	}
	return nil
}

func (dao *DAO[T]) buildSelectQuery(fields []string) string {
	if len(fields) == 0 {
		return fmt.Sprintf("SELECT * FROM %s", dao.tableName)
	}
	conditions := buildConditions(fields)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", dao.tableName, strings.Join(conditions, " AND "))
	return query
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
