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

type bulkBatch struct {
	tx   *sql.Tx
	stmt *sql.Stmt
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

// BulkInsert groups bulk items by insert statement into batches and writes them into storage
func (c *Client) BulkInsert(events []js.Object) error {
	batches := make(map[string]*bulkBatch)
	rollback := func() {
		for _, b := range batches {
			_ = b.stmt.Close()
			_ = b.tx.Rollback()
		}
	}

	items := db.ConvertEvents(events)
	for _, item := range items {
		insertStmt := insertStatement(item.Columns, item.TableName)
		batch, ok := batches[insertStmt]
		if !ok {
			tx, err := c.conn.Begin()
			if err != nil {
				rollback()
				return err
			}
			stmt, err := tx.Prepare(insertStmt)
			if err != nil {
				rollback()
				return err
			}
			batch = &bulkBatch{
				tx:   tx,
				stmt: stmt,
			}
			batches[insertStmt] = batch
		}

		_, err := batch.stmt.Exec(item.Values...)
		if err != nil {
			rollback()
			return err
		}
	}

	var err error
	for insertStmt, batch := range batches {
		delete(batches, insertStmt)
		if err = batch.tx.Commit(); err != nil {
			_ = batch.stmt.Close()
			break
		}
		if err = batch.stmt.Close(); err != nil {
			break
		}
	}
	if err != nil {
		rollback()
	}
	return err
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
