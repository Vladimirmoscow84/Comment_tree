package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (r *Router) CreateCommentHandler(c *gin.Context) {
	var req struct {
		ParentID *int   `json:"parent_id"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content cannot be empty"})
		return
	}

	comment, err := r.commentCreator.CreateComment(c.Request.Context(), req.ParentID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("[router] comment is done in DB")

	c.JSON(http.StatusCreated, comment)
}
