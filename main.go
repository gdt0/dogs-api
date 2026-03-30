package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var dataFile = "dogs.json"

type Dogs map[string][]string

func loadDogs() (Dogs, error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	var dogs Dogs
	if err := json.Unmarshal(data, &dogs); err != nil {
		return nil, err
	}
	return dogs, nil
}

func saveDogs(dogs Dogs) error {
	data, err := json.MarshalIndent(dogs, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(cors())

	// Serve static files
	r.Static("/static", "./static")

	// API routes
	api := r.Group("/api")
	{
		// GET all dogs
		api.GET("/dogs", func(c *gin.Context) {
			dogs, err := loadDogs()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, dogs)
		})

		// GET single breed
		api.GET("/dogs/:breed", func(c *gin.Context) {
			breed := strings.ToLower(c.Param("breed"))
			dogs, err := loadDogs()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if varieties, ok := dogs[breed]; ok {
				c.JSON(http.StatusOK, gin.H{"breed": breed, "varieties": varieties})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Breed not found"})
			}
		})

		// POST add breed
		api.POST("/dogs", func(c *gin.Context) {
			var input struct {
				Breed    string   `json:"breed" binding:"required"`
				Varieties []string `json:"varieties"`
			}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			breed := strings.ToLower(input.Breed)
			dogs, err := loadDogs()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if _, exists := dogs[breed]; exists {
				c.JSON(http.StatusConflict, gin.H{"error": "Breed already exists"})
				return
			}
			dogs[breed] = input.Varieties
			if err := saveDogs(dogs); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"breed": breed, "varieties": dogs[breed]})
		})

		// PUT update breed
		api.PUT("/dogs/:breed", func(c *gin.Context) {
			breed := strings.ToLower(c.Param("breed"))
			var input struct {
				Varieties []string `json:"varieties"`
			}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			dogs, err := loadDogs()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if _, exists := dogs[breed]; !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Breed not found"})
				return
			}
			dogs[breed] = input.Varieties
			if err := saveDogs(dogs); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"breed": breed, "varieties": dogs[breed]})
		})

		// DELETE breed
		api.DELETE("/dogs/:breed", func(c *gin.Context) {
			breed := strings.ToLower(c.Param("breed"))
			dogs, err := loadDogs()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if _, exists := dogs[breed]; !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Breed not found"})
				return
			}
			delete(dogs, breed)
			if err := saveDogs(dogs); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Breed '%s' deleted", breed)})
		})
	}

	// Serve index.html at root
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.File("index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	r.Run(":" + port)
}
