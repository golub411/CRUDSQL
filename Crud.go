package crudsql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents a SQLite database.
type Database struct {
	db *sql.DB
}

// OpenDatabase opens a SQLite database.
func OpenDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

// CreateTable creates a new table in the database.
func (d *Database) CreateTable(tableName string, columns []string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))
	_, err := d.db.Exec(query)
	return err
}

// InsertValue inserts a new row into a table.
func (d *Database) InsertValue(tableName string, columns []string, values []interface{}) error {
	placeholders := strings.Repeat("?, ", len(values))
	placeholders = placeholders[:len(placeholders)-2] // Remove trailing comma and space
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), placeholders)
	_, err := d.db.Exec(query, values...)
	return err
}

// SelectValue selects rows from a table.
func (d *Database) SelectValue(tableName string, columns []string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), tableName)
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				result[col] = string(b)
			} else {
				result[col] = val
			}
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

//Select value with where operator

func (d *Database) SelectValueWhere(tableName string, columns []string, where string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(columns, ", "), tableName, where)
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				result[col] = string(b)
			} else {
				result[col] = val
			}
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// UpdateValue updates rows in a table.
func (d *Database) UpdateValue(tableName string, set map[string]interface{}, where map[string]interface{}) error {
	setClause := make([]string, 0, len(set))
	whereClause := make([]string, 0, len(where))
	args := make([]interface{}, 0, len(set)+len(where))

	for col, val := range set {
		setClause = append(setClause, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	for col, val := range where {
		whereClause = append(whereClause, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(setClause, ", "), strings.Join(whereClause, " AND "))
	_, err := d.db.Exec(query, args...)
	return err
}

// DeleteValue deletes rows from a table.
func (d *Database) DeleteValue(tableName string, where map[string]interface{}) error {
	whereClause := make([]string, 0, len(where))
	args := make([]interface{}, 0, len(where))

	for col, val := range where {
		whereClause = append(whereClause, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, strings.Join(whereClause, " AND "))
	_, err := d.db.Exec(query, args...)
	return err
}
