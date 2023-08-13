package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	//router.GET("/someGet", getting)

	router.Run(":8080")
}

/*
GET /recipe_events?recipe_name=chicken%20marsala&recipe_variant=bonapetit
GET /recipes?recipe_event_id=asdf&query=foo&ingredient=onion&tag=italian
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
