package main

import (
	"recipes/models"
	storemodels "recipes/store/models"
)

func storeTagToAPI(storeTag *storemodels.Tag) models.Tag {
	return models.Tag(storeTag.Name)
}
