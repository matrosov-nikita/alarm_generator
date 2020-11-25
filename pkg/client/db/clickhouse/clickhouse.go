package clickhouse

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db"

	"github.com/ClickHouse/clickhouse-go"
	js "github.com/itimofeev/go-util/json"
)

type Client struct {
	conn *sql.DB
}

func NewClient(url string) (*Client, error) {
	conn, err := sql.Open("clickhouse", url)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) BulkInsert(events []js.Object) error {
	if len(events) == 0 {
		return nil
	}

	items := db.ConvertEvents(events)

	tx, err := c.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertStatement(items[0].Columns, items[0].TableName))
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("statement close failed: %v\n", err)
		}
	}()

	for _, item := range items {
		if _, err := stmt.Exec(item.Values...); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func insertStatement(columns []string, tableName string) string {
	/* #nosec */
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ","),
		placeholdersString(len(columns)))
}

func placeholdersString(count int) string {
	if count == 0 {
		return ""
	}
	placeholders := strings.Repeat("?,", count)
	return placeholders[:len(placeholders)-1]
}
