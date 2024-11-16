package db

import (
	"fmt"
)

type BaseDAO[T any] struct {
	DB        *DB
	tableName string
}

func newBaseDao[T any](db *DB, tableName string) *BaseDAO[T] {
	return &BaseDAO[T]{
		DB:        db,
		tableName: tableName,
	}
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
