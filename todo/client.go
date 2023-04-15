package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is a Go HTTP client to interact with the Todo API.
type Client struct {
	baseURL *url.URL
	http    *http.Client
}

// NewClient creates a new Client using rawURL as the base URL for the Todo
// API.
func NewClient(rawURL string) (*Client, error) {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	c := Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}

	return &c, nil
}

// ListTodos retrieves a list of all todos from the API.
func (c *Client) ListTodos() ([]Todo, error) {
	u := c.baseURL.JoinPath("/api/todo")

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %v", resp.StatusCode)
	}

	todos := make([]Todo, 0)
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, err
	}

	return todos, nil
}

// GetTodo retrieves a single todo by its id from the API.
func (c *Client) GetTodo(id string) (Todo, error) {
	u := c.baseURL.JoinPath("/api/todo", id)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return Todo{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Todo{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Todo{}, fmt.Errorf("invalid status code %v", resp.StatusCode)
	}

	var td Todo
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		return Todo{}, err
	}

	return td, nil
}

// CreateTodo creates a todo.
func (c *Client) CreateTodo(params TodoCreateParams) (Todo, error) {
	u := c.baseURL.JoinPath("/api/todo")

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(params); err != nil {
		return Todo{}, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return Todo{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Todo{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return Todo{}, fmt.Errorf("invalid status code %v", resp.StatusCode)
	}

	var td Todo
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		return Todo{}, err
	}

	return td, nil
}

// UpdateTodo updates an existing todo given by id.
func (c *Client) UpdateTodo(id string, params TodoUpdateParams) (Todo, error) {
	u := c.baseURL.JoinPath("/api/todo", id)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(params); err != nil {
		return Todo{}, err
	}

	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return Todo{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Todo{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Todo{}, fmt.Errorf("invalid status code %v", resp.StatusCode)
	}

	var td Todo
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		return Todo{}, err
	}

	return td, nil
}

// DeleteTodo deletes a todo by its id.
func (c *Client) DeleteTodo(id string) error {
	u := c.baseURL.JoinPath("/api/todo", id)

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid status code %v", resp.StatusCode)
	}

	return nil
}
