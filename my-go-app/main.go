package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

func ListRecipeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"recipes": recipes,
	})
}

func IndexHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.JSON(200, gin.H{
		"title": name, "subtitle": "Greetings, ",
		"message": "Welcome to the Index Page!", "version": "1.0",
		"date": "2022-10-20", "author": "John Doe",
		"description": "This is a sample API response", "status": "success",
	})
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	// c.JSON(
	// 	http.StatusOK,
	// 	recipe,
	// )
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Recipe added successfully", "recipe": recipe})
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var updatedRecipe Recipe
	if err := c.ShouldBindJSON(&updatedRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, r := range recipes {
		if r.ID == id {
			recipes[i] = updatedRecipe
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Recipe updated successfully", "recipe": updatedRecipe})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for i, r := range recipes {
		if r.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	var deletedRecipeName string
	if index != -1 {
		// Remove the recipe from the slice
		deletedRecipeName = recipes[index].Name
		recipes = append(recipes[:index], recipes[index+1:]...)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "recipe has been deleted", "deletedRecipeName": deletedRecipeName})
}

func SearchRecipeHandler(c *gin.Context) {
	query := c.Query("tag")
	results := make([]Recipe, 0)
	for _, recipe := range recipes {
		for _, tag := range recipe.Tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(query)) {
				results = append(results, recipe)
				break
			}
		}
	}
	c.JSON(200, gin.H{"results": results})
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run(":8000")
}
