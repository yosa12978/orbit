package dto

import "time"

type SnippetResponse struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type SnippetCreateRequest struct {
	Content string `json:"content"`
}

type SnippetCreateResponse struct {
	ID int64 `json:"id"`
}
