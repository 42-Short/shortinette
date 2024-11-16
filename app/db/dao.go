package db

import (
	"fmt"
	"reflect"
	"strings"
)

type BaseDAO[T any] struct {
	DB        *DB
	tableName string
	tags      []string
}

func (dao *BaseDAO[T]) Insert(data *T) error {
	query, err := dao.buildInsertQuery()
	if err != nil {
		return fmt.Errorf("failed to build insertion query")
	}
	_, err = dao.DB.namedExecWithTimeout(query, data)
	if err != nil {
		return fmt.Errorf("failed to insert data into table %s: %v", dao.tableName, err)
	}

	return nil
}

func (dao *BaseDAO[T]) GetAll() ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s;", dao.tableName)

	var retrievedData []T
	err := dao.DB.selectWithTimeout(&retrievedData, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}

	return retrievedData, nil
}

func newBaseDao[T any](db *DB, tableName string) *BaseDAO[T] {
	var dummy T
	tags := extractStructTags(&dummy)
	return &BaseDAO[T]{
		DB:        db,
		tableName: tableName,
		tags:      tags,
	}
}

func (dao *BaseDAO[T]) buildInsertQuery() (string, error) {
	columns := strings.Join(dao.tags, ", ")
	placeholders := strings.Join(createNamedPlaceholders(dao.tags), ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", dao.tableName, columns, placeholders)

	return query, nil
}

func createNamedPlaceholders(tags []string) []string {
	placeholders := make([]string, len(tags))
	for i, tag := range tags {
		placeholders[i] = ":" + tag
	}
	return placeholders
}

func extractStructTags[T any](data *T) []string {
	t := reflect.TypeOf(*data)

	dbTags := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		dbTags = append(dbTags, dbTag)
	}
	return dbTags
}
