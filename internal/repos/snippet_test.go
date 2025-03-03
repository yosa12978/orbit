package repos_test

import (
	"context"
	"fmt"
	"log/slog"
	"orbit-app/internal/repos"
	"orbit-app/internal/snippet"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func getRepo(t *testing.T) snippet.Repo {
	miniInstance := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: miniInstance.Addr()})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return repos.NewSnippetRepoRedis(rdb, logger)
}

func seedRepo(ctx context.Context, repo snippet.Repo) error {
	for i := 0; i < 10; i++ {
		id, err := snippet.NewID(int64(i))
		if err != nil {
			return err
		}
		newSnippet, err := snippet.NewFromPrimitives(
			int64(id),
			fmt.Sprintf("content #%d", i),
			time.Now().UTC(),
		)
		if err != nil {
			return err
		}
		_, err = repo.Create(ctx, newSnippet)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCreateAndGetSnippetByID(t *testing.T) {
	repo := getRepo(t)
	if err := seedRepo(t.Context(), repo); err != nil {
		t.Errorf("error while seeding db: %s", err.Error())
		return
	}
	for i := 0; i < 10; i++ {
		requestedID, err := snippet.NewID(int64(i))
		if err != nil {
			t.Error(err.Error())
			return
		}
		snip, err := repo.GetByID(t.Context(), requestedID)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if snip.ID() != requestedID {
			t.Error("snip.ID() diesn't match expected id")
			return
		}
		t.Logf("%+v\n", snip)
	}
}

func TestUpdateSnippet(t *testing.T) {
	repo := getRepo(t)
	if err := seedRepo(t.Context(), repo); err != nil {
		t.Errorf("error while seeding db: %s", err.Error())
	}
	for i := 0; i < 10; i++ {
		requestedID, err := snippet.NewID(int64(i))
		if err != nil {
			t.Error(err.Error())
			return
		}
		snip, err := repo.GetByID(t.Context(), requestedID)
		updatedSnip, _ := snippet.New(snip.ID(), "updated content", snip.CreatedAt())
		snipID, err := repo.Update(t.Context(), &updatedSnip)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if snipID != updatedSnip.ID() {
			t.Error("id's don't match")
			return
		}
		//checking
		snipp, err := repo.GetByID(t.Context(), requestedID)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if snipp.Content() != "updated content" {
			t.Error("content doesn't match")
			return
		}
	}
}

func TestDeleteSnippet(t *testing.T) {
	repo := getRepo(t)
	if err := seedRepo(t.Context(), repo); err != nil {
		t.Errorf("error while seeding db: %s", err.Error())
		return
	}
	reqID, err := snippet.NewID(0)
	if err != nil {
		t.Error(err.Error())
		return
	}
	id, err := repo.Delete(t.Context(), reqID)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if reqID != id {
		t.Error("id's doesn't match")
		return
	}
	if _, err := repo.GetByID(t.Context(), id); err == nil {
		t.Error("it still exist")
		return
	}
}
