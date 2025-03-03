package snippet

type CacheModel struct {
	ID        ID        `json:"id"`
	Content   Content   `json:"content"`
	CreatedAt CreatedAt `json:"createdAt"`
}

func NewCacheModel(s Snippet) CacheModel {
	return CacheModel{
		ID:        s.ID(),
		Content:   s.Content(),
		CreatedAt: s.CreatedAt(),
	}
}

func (c CacheModel) Domain() Snippet {
	return Snippet{
		id:        c.ID,
		content:   c.Content,
		createdAt: c.CreatedAt,
	}
}
