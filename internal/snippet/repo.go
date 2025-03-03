package snippet

import (
	"context"
)

type Repo interface {
	GetByID(ctx context.Context, id ID) (Snippet, error)
	Create(ctx context.Context, snip Snippet) (ID, error)
	Update(ctx context.Context, snip *Snippet) (ID, error)
	Delete(ctx context.Context, id ID) (ID, error)
}
