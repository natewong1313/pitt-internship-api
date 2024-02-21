package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natewong1313/pitt-internship-api/scraper"
	gocache "github.com/patrickmn/go-cache"
)

func main() {
	// Cache the job listings for 1 hour
	cache := gocache.New(1*time.Hour, 10*time.Minute)
	cache.Set("listings", scraper.Scrape(), gocache.DefaultExpiration)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		var listings []scraper.JobListing
		cachedListings, found := cache.Get("listings")
		if !found {
			listings = scraper.Scrape()
			cache.Set("listings", listings, gocache.DefaultExpiration)
		} else {
			listings = cachedListings.([]scraper.JobListing)
		}
		c.JSON(200, gin.H{
			"listings": listings,
		})
	})
	router.Run()
}
