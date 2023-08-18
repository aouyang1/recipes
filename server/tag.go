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
	PostTagTimeout   = time.Duration(5 * time.Second)
	DeleteTagTimeout = time.Duration(5 * time.Second)
)

// PostTag adds a tag to the db
func (s *Server) PostTag(c *gin.Context) {
	var req models.Tag
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), PostTagTimeout)
	defer cancel()
	tagID, err := s.postTag(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	req.ID = tagID
	c.JSON(http.StatusCreated, req)
}

func (s *Server) postTag(ctx context.Context, tag *models.Tag) (uint64, error) {
	t := &storemodels.Tag{
		ID:   tag.ID,
		Name: tag.Name,
	}
	return s.store.UpsertTag(ctx, t)
}

// DeleteTag removes a tag from the db
func (s *Server) DeleteTag(c *gin.Context) {
	var req models.Tag
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), DeleteTagTimeout)
	defer cancel()
	if err := s.deleteTag(ctx, &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) deleteTag(ctx context.Context, tag *models.Tag) error {
	return s.store.DeleteTag(ctx, tag.ID)
}
