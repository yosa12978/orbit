package services

import (
	"context"
	"log/slog"
	"orbit-app/internal/config"
	"orbit-app/internal/dto"
	"orbit-app/internal/snippet"
	"orbit-app/pkg"
	"time"
)

type SnippetService interface {
	GetByID(ctx context.Context, id int64) (dto.SnippetResponse, error)
	Create(ctx context.Context, snipDto dto.SnippetCreateRequest) (dto.SnippetCreateResponse, error)
}

type snippetService struct {
	repo   snippet.Repo
	logger *slog.Logger
}

func NewSnippetService(repo snippet.Repo, logger *slog.Logger) SnippetService {
	return &snippetService{
		repo:   repo,
		logger: logger,
	}
}

func (s *snippetService) Create(
	ctx context.Context,
	snipDto dto.SnippetCreateRequest,
) (dto.SnippetCreateResponse, error) {
	id := pkg.NewSfGen(config.Get().Server.InstanceID).ID()
	snip, err := snippet.NewFromPrimitives(
		id,
		snipDto.Content,
		time.Now().UTC(),
	)
	if err != nil {
		return dto.SnippetCreateResponse{}, err
	}
	_, err = s.repo.Create(ctx, snip)
	return dto.SnippetCreateResponse{ID: id}, nil
}

func (s *snippetService) GetByID(
	ctx context.Context,
	id int64,
) (dto.SnippetResponse, error) {
	snipID, err := snippet.NewID(id)
	if err != nil {
		return dto.SnippetResponse{}, err
	}
	snip, err := s.repo.GetByID(ctx, snipID)
	return dto.SnippetResponse{
		ID:        int64(snip.ID()),
		Content:   string(snip.Content()),
		CreatedAt: time.Time(snip.CreatedAt()),
	}, err
}
