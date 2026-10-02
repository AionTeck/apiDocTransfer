package postman_client

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.getpostman.com"

type PostmanClient struct {
	APIKey      string
	WorkspaceID string
	HTTPClient  *http.Client
}

func NewPostmanClient(apiKey string, WorkspaceID string, client *http.Client) *PostmanClient {
	return &PostmanClient{
		APIKey:      apiKey,
		WorkspaceID: WorkspaceID,
		HTTPClient:  client,
	}
}

func (c *PostmanClient) UpsertFolder(ctx context.Context, targetID, folderName string, converted map[string]any) error {
	var got struct {
		Collection map[string]any `json:"collection"`
	}
	if err := c.do(ctx, http.MethodGet, "/collections/"+targetID, nil, &got); err != nil {
		return fmt.Errorf("get target collection: %w", err)
	}
	col := got.Collection

	folder := map[string]any{
		"name": folderName,
		"item": converted["item"],
	}
	if auth, ok := converted["auth"]; ok {
		folder["auth"] = auth
	}

	items, _ := col["item"].([]any)
	replaced := false
	for i, it := range items {
		if m, ok := it.(map[string]any); ok && m["name"] == folderName {
			items[i] = folder
			replaced = true
			break
		}
	}
	if !replaced {
		items = append(items, folder)
	}
	col["item"] = items

	if err := c.do(ctx, http.MethodPut, "/collections/"+targetID,
		map[string]any{"collection": col}, nil); err != nil {
		return fmt.Errorf("put target collection: %w", err)
	}
	return nil
}

func (c *PostmanClient) do(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: %s: %s", method, path, resp.Status, data)
	}
	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}
