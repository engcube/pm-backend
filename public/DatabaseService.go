package PrivateMessageBackendPublic

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFILE = "pmbackend.sqlite3"
)

var (
	DBConnection *sql.DB
)

// DB 获取DB连接
func DB() (*sql.DB, error) {
	var err error
	if DBConnection == nil {
		DBConnection, err = sql.Open("sqlite3", DBFILE)
		if err != nil {
			return nil, err
		}
	}
	return DBConnection, nil
}

// Insert 插入操作
func Insert(sql string, args ...interface{}) (int64, error) {
	db, err := DB()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update 更新操作(支持update,delete)
func Update(sql string, args ...interface{}) (int64, error) {
	db, err := DB()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

// Select 选择操作
func Select(sqlQuery string, args ...interface{}) ([][]string, error) {
	db, err := DB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	res := make([][]string, 0)
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		line := make([]string, 0)
		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			line = append(line, value)
		}
		res = append(res, line)
	}
	return res, nil
}
