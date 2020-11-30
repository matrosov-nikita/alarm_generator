package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	js "github.com/itimofeev/go-util/json"
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) BulkInsert(items []js.Object) error {
	body, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to encode batch events: %+v", err)
	}
	log.Println("items", string(body))

	response, err := http.Post(c.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send events to %s: %v", c.url, err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send events, http error code=%d", response.StatusCode)
	}

	return nil
}
