package store

import (
	"context"
	"recipes/store/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertRecipeEventContext(t *testing.T) {
	testData := map[string]struct {
		recipeEvent []*models.RecipeEvent
		err         error
	}{
		"nil event": {
			[]*models.RecipeEvent{nil}, ErrNilRecipeEvent,
		},
		"insert new event": {
			[]*models.RecipeEvent{
				{
					ID:           "asdf",
					ScheduleDate: 1234,
					Title:        "myrecipe",
					Description:  "Step1: bake",
				},
			},
			nil,
		},
		"update event": {
			[]*models.RecipeEvent{
				{
					ID:           "asdf",
					ScheduleDate: 1234,
					Title:        "myrecipe",
					Description:  "Step1: bake",
				},
				{
					ID:           "asdf",
					ScheduleDate: 4321,
					Title:        "myrecipe2",
					Description:  "Step2: cook",
				},
			},
			nil,
		},
	}

	client, err := NewTestClient()
	require.Nil(t, err)

	for name, td := range testData {
		t.Run(name, func(t *testing.T) {
			defer Cleanup(client.conn)

			for _, revent := range td.recipeEvent {
				if err := client.UpsertRecipeEvent(context.Background(), revent); err != nil {
					assert.ErrorIs(t, err, td.err)
					return
				}
				assert.Nil(t, td.err)
			}

			exists, err := client.ExistsRecipeEvent(context.Background(), td.recipeEvent[len(td.recipeEvent)-1])
			require.Nil(t, err)
			assert.True(t, exists)
		})
	}
}
