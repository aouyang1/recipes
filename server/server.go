package main

import (
	"fmt"
	"os"

	"recipes/store"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Server struct {
	router *gin.Engine
	store  *store.Client
}

func NewServer() (*Server, error) {
	storeClient, err := store.NewClient(
		&mysql.Config{
			User:   os.Getenv("USER_MYSQL_USERNAME"),
			Passwd: os.Getenv("USER_MYSQL_PASSWORD"),
			DBName: os.Getenv("USER_MYSQL_DB"),
			Net:    "tcp",
			Addr:   "localhost",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create store client, %w", err)
	}

	router := gin.Default()

	s := &Server{
		router: router,
		store:  storeClient,
	}

	router.GET("/recipe_events", s.GetRecipeEvents)
	router.GET("/recipes", s.GetRecipes)
	router.POST("/recipe", s.PostRecipe)
	router.POST("/recipe_tag", s.PostRecipeTag)
	router.DELETE("/recipe_tag", s.DeleteRecipeTag)
	router.POST("/recipe_ingredient", s.PostRecipeIngredient)
	router.DELETE("/recipe_ingredient", s.DeleteRecipeIngredient)

	return s, nil
}

func (s *Server) Run() error {
	if err := s.router.Run(":8080"); err != nil {
		return err
	}
	return nil
}

/*

GET /tags
POST /tag {name: italian}
DELETE /tag {name: italian}

GET /ingredients
POST /ingredient {name: onion}
DELETE /ingredient {name: onion}
*/
