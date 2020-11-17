package postgres

import (
	"database/sql"
	"log"

	"github.com/matrosov-nikita/alarm_generator/event"

	"github.com/lib/pq"
)

type Client struct {
	conn *sql.DB
}

type bulkBatch struct {
	tx   *sql.Tx
	stmt *sql.Stmt
}

func NewClient(url string) (*Client, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		if postgresError, ok := err.(*pq.Error); ok {
			log.Println(postgresError.Code, postgresError.Message)
		}
		return nil, err
	}

	return &Client{conn: conn}, nil
}

// Close closes connection to Postgres
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnection returns connection to Postgres
func (c *Client) GetConnection() *sql.DB {
	return c.conn
}

// BulkInsert groups bulk items by insert statement into batches and writes them into storage
func (c *Client) BulkInsert(items []event.Item) error {
	batches := make(map[string]*bulkBatch)
	rollback := func() {
		for _, b := range batches {
			_ = b.stmt.Close()
			_ = b.tx.Rollback()
		}
	}
	for _, item := range items {
		insertStmt := item.InsertStatement(item.Columns())
		batch, ok := batches[insertStmt]
		if !ok {
			tx, err := c.conn.Begin()
			if err != nil {
				rollback()
				return err
			}
			stmt, err := tx.Prepare(pq.CopyIn(item.TableName(), item.Columns()...))
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
		_, err := batch.stmt.Exec(item.Values()...)
		if err != nil {
			rollback()
			return err
		}
	}

	var err error
	for insertStmt, batch := range batches {
		delete(batches, insertStmt)
		if _, err = batch.stmt.Exec(); err != nil {
			break
		}
		if err = batch.stmt.Close(); err != nil {
			break
		}

		if err = batch.tx.Commit(); err != nil {
			break
		}
	}
	if err != nil {
		rollback()
	}
	return err
}

// RunInTransaction runs specified function with transaction context
func (c *Client) RunInTransaction(fn func(*sql.Tx) error) error {
	tx, err := c.conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println("Postgres transaction rollback error")
			}
			panic(err)
		}
	}()
	if err := fn(tx); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Println("Postgres transaction rollback error")
		}
		return err
	}
	return tx.Commit()
}

// GetName returns storage name
func (c *Client) GetName() string {
	return "postgres"
}
