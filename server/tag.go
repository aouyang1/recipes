package main

import (
	"context"
	"net/http"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	PostTagTimeout = time.Duration(5 * time.Second)
)

type TagRequest struct {
	Name models.Tag `json:"name"`
}

func (s *Server) PostTag(c *gin.Context) {
	var req TagRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), PostTagTimeout)
	defer cancel()
	if err := s.postTag(ctx, req.Name); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) postTag(ctx context.Context, tag models.Tag) error {
	t := &storemodels.Tag{
		ID:   tag.ID(),
		Name: string(tag),
	}
	return s.store.InsertTag(ctx, t)
}
