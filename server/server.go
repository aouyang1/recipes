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

	router.LoadHTMLFiles("static/index.html")
	router.Static("/static", "static")

	router.GET("/", s.Index)
	router.GET("/recipe_events", s.GetRecipeEvents)

	router.GET("/recipes", s.GetRecipes)
	router.POST("/recipe", s.PostRecipe)
	router.PUT("/recipe", s.UpdateRecipe)

	router.GET("/tags", s.GetTags)
	router.POST("/tag", s.PostTag)
	router.DELETE("/tag", s.DeleteTag)

	router.GET("/ingredients", s.GetIngredients)
	router.POST("/ingredient", s.PostIngredient)
	router.DELETE("/ingredient", s.DeleteIngredient)

	return s, nil
}

func (s *Server) Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func (s *Server) Run() error {
	if err := s.router.Run(":8080"); err != nil {
		return err
	}
	return nil
}
