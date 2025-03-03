package snippet

import "time"

type ID int64
type Content string
type CreatedAt time.Time

func NewID(id int64) (ID, error) {
	return ID(id), nil
}

func NewContent(content string) (Content, error) {
	if err := validateContent(content); err != nil {
		return Content(""), err
	}
	return Content(content), nil
}

func NewCreatedAt(createdAt time.Time) (CreatedAt, error) {
	return CreatedAt(createdAt), nil
}
