package repos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"orbit-app/internal/snippet"
	"time"

	"github.com/redis/go-redis/v9"
)

type snippetRepoRedis struct {
	rdb    redis.Cmdable
	logger *slog.Logger
}

func NewSnippetRepoRedis(rdb redis.Cmdable, logger *slog.Logger) snippet.Repo {
	return &snippetRepoRedis{
		rdb:    rdb,
		logger: logger,
	}
}

func (s *snippetRepoRedis) getDefaultSnippetTTL() time.Duration {
	return 7 * 24 * time.Hour
}

func (s *snippetRepoRedis) GetByID(
	ctx context.Context,
	id snippet.ID,
) (snippet.Snippet, error) {
	var snip snippet.CacheModel
	key := fmt.Sprintf("snippets:%v", id)
	res, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return snippet.Snippet{}, err
	}
	if err := json.Unmarshal([]byte(res), &snip); err != nil {
		return snippet.Snippet{}, err
	}
	return snip.Domain(), nil
}

func (s *snippetRepoRedis) Create(
	ctx context.Context,
	snip snippet.Snippet,
) (snippet.ID, error) {
	key := fmt.Sprintf("snippets:%v", snip.ID())
	cacheModel := snippet.NewCacheModel(snip)
	snipJSON, err := json.Marshal(cacheModel)
	if err != nil {
		return snippet.ID(0), err
	}
	_, err = s.rdb.Set(ctx, key, snipJSON, s.getDefaultSnippetTTL()).Result()
	return snip.ID(), err
}

func (s *snippetRepoRedis) Update(
	ctx context.Context,
	snip *snippet.Snippet,
) (snippet.ID, error) {
	if snip == nil {
		return snippet.ID(0), errors.New("nil snippet")
	}

	id := (*snip).ID()
	key := fmt.Sprintf("snippets:%v", id)

	ttl, _ := s.rdb.TTL(ctx, key).Result()
	if ttl == -2 {
		return snippet.ID(0), errors.New("snippet doesn't exist")
	}

	cacheModel := snippet.NewCacheModel(*snip)
	cacheModelJSON, _ := json.Marshal(cacheModel)
	return id, s.rdb.Set(ctx, key, cacheModelJSON, ttl).Err()
}

func (s *snippetRepoRedis) Delete(
	ctx context.Context,
	id snippet.ID,
) (snippet.ID, error) {
	key := fmt.Sprintf("snippets:%v", id)
	return id, s.rdb.Del(ctx, key).Err()
}
