package handlers

import (
	"comment_tree/internal/model"
	"context"

	"github.com/gin-gonic/gin"
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

func New(router *ginext.Engine, creator commentCreator, getter commentGetter, deleter commentDeleter) *Router {
	return &Router{
		Router:         router,
		commentCreator: creator,
		commentGetter:  getter,
		commentDeleter: deleter,
	}
}

func (r *Router) Routes() {
	r.Router.POST("/comments", r.CreateCommentHandler)
	r.Router.GET("/comments", r.GetCommentHandler)
	r.Router.DELETE("/comments/:id", r.DeleteCommentHandler)
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")
}
