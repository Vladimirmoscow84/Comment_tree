package service

import (
	"comment_tree/internal/model"
	"context"
)

type commentGetterRepo interface {
	CreateComment(ctx context.Context, parentID *int, content string) (*model.Comment, error)
	GetSubtree(ctx context.Context, id int) ([]model.Comment, error)
	DeleteSubtree(ctx context.Context, id int) error
	SearchComments(ctx context.Context, queryText string, limit, offset int) ([]model.Comment, error)
	ListRootComments(ctx context.Context, limit, offset int) ([]model.Comment, error)
}

type Service struct {
	storage commentGetterRepo
}

func New(storage commentGetterRepo) *Service {
	return &Service{
		storage: storage,
	}
}
