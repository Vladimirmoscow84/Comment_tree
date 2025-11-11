package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) GetCommentHandler(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Request.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	if search := query.Get("search"); search != "" {
		results, err := r.commentGetter.SearchComments(ctx, search, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, results)
		return
	}

	if parent := query.Get("parent"); parent != "" {
		id, err := strconv.Atoi(parent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent id"})
			return
		}
		tree, err := r.commentGetter.GetSubtree(ctx, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tree)
		return
	}

	comments, err := r.commentGetter.ListRootComments(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
