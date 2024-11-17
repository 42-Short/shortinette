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

func NewDAO[T any](db *DB, tableName string) *DAO[T] {
	var dummy T
	tags := extractTags(&dummy, "db")
	if len(tags) == 0 {
		panic("NewBaseDao: Expected database tags (db) but got 0")
	}
	primaryKeys := extractTags(&dummy, "primaryKey")
	if len(primaryKeys) == 0 {
		panic("NewBaseDao: Expected primaryKey tags (primaryKey) but got 0")
	}
	return &DAO[T]{
		DB:          db,
		tableName:   tableName,
		dbTags:      tags,
		primaryKeys: primaryKeys,
	}
}

func (dao *DAO[T]) Insert(data *T) error {
	query := dao.buildInsertQuery()
	_, err := dao.DB.namedExecWithTimeout(query, data)
	if err != nil {
		return fmt.Errorf("failed to insert data into table %s: %v", dao.tableName, err)
	}
	return nil
}

func (dao *DAO[T]) GetAll() ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s;", dao.tableName)

	var retrievedData []T
	err := dao.DB.selectWithTimeout(&retrievedData, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}
	return retrievedData, nil
}

func (dao *DAO[T]) Get(args ...any) (*T, error) {
	var retrievedData T

	query := dao.buildGetQuery(dao.primaryKeys)
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

	query := dao.buildGetQuery(fields)

	var results []T
	err := dao.DB.selectWithTimeout(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter data: %v", err)
	}

	return results, nil
}

func (dao *DAO[T]) buildGetQuery(fields []string) string {
	conditions := make([]string, 0, len(fields))
	for _, field := range fields {
		conditions = append(conditions, fmt.Sprintf("%s = ?", field))
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", dao.tableName, strings.Join(conditions, " AND "))
	return query
}

func (dao *DAO[T]) buildInsertQuery() string {
	columns := strings.Join(dao.dbTags, ", ")
	placeholders := strings.Join(createNamedPlaceholders(dao.dbTags), ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", dao.tableName, columns, placeholders)
	return query
}

func createNamedPlaceholders(tags []string) []string {
	placeholders := make([]string, 0, len(tags))
	for _, tag := range tags {
		placeholders = append(placeholders, ":"+tag)
	}
	return placeholders
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
