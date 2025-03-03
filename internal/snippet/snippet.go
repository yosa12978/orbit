package snippet

import (
	"orbit-app/pkg"
	"time"
)

type Snippet struct {
	id        ID
	content   Content
	createdAt CreatedAt
}

func New(id ID, content Content, createdAt CreatedAt) (Snippet, error) {
	return Snippet{
		id:        id,
		content:   content,
		createdAt: createdAt,
	}, nil
}

func NewFromPrimitives(id int64, content string, createdAt time.Time) (Snippet, error) {
	problems := pkg.ValidationError{}
	newID, err := NewID(id)
	if err != nil {
		problems["id"] = err.Error()
	}
	newContent, err := NewContent(content)
	if err != nil {
		problems["content"] = err.Error()
	}
	newCreatedAt, err := NewCreatedAt(createdAt)
	if err != nil {
		problems["createdAt"] = err.Error()
	}
	if len(problems) > 0 {
		return Snippet{}, problems
	}
	snippet := Snippet{
		id:        newID,
		content:   newContent,
		createdAt: newCreatedAt,
	}
	return snippet, nil
}

func (s Snippet) ID() ID {
	return s.id
}

func (s Snippet) Content() Content {
	return s.content
}

func (s Snippet) CreatedAt() CreatedAt {
	return s.createdAt
}
