package main

import (
	"github.com/gin-gonic/gin"
	"github.com/natewong1313/pitt-intern-api/scraper"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"listings": scraper.Scrape(),
		})
	})
	r.Run()
}
