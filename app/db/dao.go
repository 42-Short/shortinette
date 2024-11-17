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

func NewBaseDao[T any](db *DB, tableName string) *BaseDAO[T] {
	var dummy T
	tags := extractStructTags(&dummy)
	if len(tags) == 0 {
		panic("NewBaseDao: Expected database tags tags but got none")
	}
	return &BaseDAO[T]{
		DB:        db,
		tableName: tableName,
		tags:      tags,
	}
}

func (dao *BaseDAO[T]) Insert(data *T) error {
	query := dao.buildInsertQuery()
	_, err := dao.DB.namedExecWithTimeout(query, data)
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

func (dao *BaseDAO[T]) Get(columnNames []string, args ...any) (*T, error) {
	var retrievedData T
	query := dao.buildGetQuery(columnNames)
	err := dao.DB.getWithTimeout(&retrievedData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from table %s: %v", dao.tableName, err)
	}
	return &retrievedData, err
}

func (dao *BaseDAO[T]) buildGetQuery(columnNames []string) string {
	conditions := make([]string, 0, len(columnNames))
	for _, columnName := range columnNames {
		conditions = append(conditions, fmt.Sprintf("%s = ?", columnName))
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", dao.tableName, strings.Join(conditions, " AND "))
	return query
}

func (dao *BaseDAO[T]) buildInsertQuery() string {
	columns := strings.Join(dao.tags, ", ")
	placeholders := strings.Join(createNamedPlaceholders(dao.tags), ", ")
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

func extractStructTags[T any](data *T) []string {
	t := reflect.TypeOf(*data)

	dbTags := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		dbTags = append(dbTags, dbTag)
	}
	return dbTags
}
