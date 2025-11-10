package handlers

import (
	"comment_tree/internal/model"
	"context"

	"github.com/!vladimirmoscow84/wbf/ginext"
	"github.com/wb-go/wbf/ginext"
)

type commentCreator interface {
	CreateComment(ctx context.Context, parentID *int, content string) (*model.Comment, error)
}

type commentGetter interface {
	GetSubtree(ctx context.Context, id int) ([]model.Comment, error)
	SearchComments(ctx context.Context, queryText string, limit, offset int) ([]model.Comment, error)
	ListRootComments(ctx context.Context, limit, offset int) ([]model.Comment, error)
}

type commentDeleter interface {
	DeleteSubtree(ctx context.Context, id int) error
}

type Router struct {
	Router         *ginext.Engine
	commentCreator commentCreator
	commentGetter  commentGetter
	commentDeleter commentDeleter
}
