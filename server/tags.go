package main

import (
	"context"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	GetTagsTimeout = time.Duration(5 * time.Second)
)

// GetTags returns all tags
func (s *Server) GetTags(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), GetTagsTimeout)
	defer cancel()
	tags, err := s.getTags(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tags)

}

func (s *Server) getTags(ctx context.Context) ([]*models.Tag, error) {
	storeTags, err := s.store.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]*models.Tag, 0, len(storeTags))
	for _, storeTag := range storeTags {
		tags = append(tags, storeTagToAPI(storeTag))
	}

	return tags, nil
}

func storeTagToAPI(storeTag *storemodels.Tag) *models.Tag {
	return &models.Tag{
		ID:    storeTag.ID,
		Name:  storeTag.Name,
		Count: storeTag.Count,
	}
}
