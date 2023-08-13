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

	return s, nil
}

func (s *Server) Run() error {
	if err := s.router.Run(":8080"); err != nil {
		return err
	}
	return nil
}

/*
GET /recipe?name=chicken%20marsala&variant=bonapetit
POST /recipe?name=chicken%20marsala&variant=bonapetit

POST /recipe_tag {recipe_name: chicken marsala, recipe_variant: bonapetit, tag: italian}
DELETE /recipe_tag {recipe_name: chicken marsala, recipe_variant: bonapetit, tag: italian}

GET /tags
POST /tag {name: italian}
DELETE /tag {name: italian}

POST /recipe_ingredient {recipe_name: "chicken marsala", recipe_variant: "bonapetit", ingredient: "onion", quantity: 1, unit: "cup", size: ""}
DELETE /recipe_ingredient {recipe_name: "chicken marsala", recipe_variant: "bonapetit", ingredient: "onion"}

GET /ingredients
POST /ingredient {name: onion}
DELETE /ingredient {name: onion}
*/