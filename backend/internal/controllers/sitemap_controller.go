package controllers

import (
	"encoding/xml"
	"net/http"
	"os"
	"time"

	"website-dummy/internal/config"
	"website-dummy/internal/models"

	"github.com/gin-gonic/gin"
)

type SitemapController struct{}

// NewSitemapController creates a new sitemap controller
func NewSitemapController() *SitemapController {
	return &SitemapController{}
}

// URL represents a URL entry in the sitemap
type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	LastMod    string   `xml:"lastmod"`
	ChangeFreq string   `xml:"changefreq"`
	Priority   string   `xml:"priority"`
}

// Sitemap represents the sitemap structure
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// GetSitemap generates and returns a sitemap.xml
func (sc *SitemapController) GetSitemap(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	// Get base URL from environment or use default
	baseURL := os.Getenv("SITE_URL")
	if baseURL == "" {
		baseURL = "https://website-dummy.com"
	}
	// Remove trailing slash if present
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	urls := []URL{}

	// Static pages
	staticPages := []struct {
		path        string
		changeFreq string
		priority    string
	}{
		{"/", "daily", "1.0"},
		{"/about", "monthly", "0.8"},
		{"/search", "daily", "0.9"},
		{"/buy", "monthly", "0.8"},
		{"/sell", "monthly", "0.8"},
		{"/contact", "monthly", "0.7"},
		{"/rentals", "monthly", "0.8"},
		{"/blog", "weekly", "0.8"},
	}

	for _, page := range staticPages {
		urls = append(urls, URL{
			Loc:        baseURL + page.path,
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: page.changeFreq,
			Priority:   page.priority,
		})
	}

	// Team member pages
	teamMembers := []struct {
		slug string
	}{
		{"julie-pogue"},
		{"guthrie-zaring"},
		{"alissa-meriwether"},
		{"amanda-webb"},
		{"erin-blascak"},
		{"kevin-grebe"},
		{"dylan-hogan"},
		{"cameron-renneisen"},
		{"lesa-seibert"},
		{"sheila-weber"},
	}

	for _, member := range teamMembers {
		urls = append(urls, URL{
			Loc:        baseURL + "/agent/" + member.slug,
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   "0.6",
		})
	}

	// Neighborhood pages
	neighborhoods := []string{
		"anchorage-homes-for-sale",
		"audubon-park-homes-for-sale",
		"cherokee-triangle-homes-for-sale",
		"crescent-hill-homes-for-sale",
		"forest-springs-homes-for-sale",
		"highlands-homes-for-sale",
		"hunting-creek-homes-for-sale",
		"indian-hills-homes-for-sale",
		"lake-forest-homes-for-sale",
		"mockingbird-valley-homes-for-sale",
		"norton-commons-homes-for-sale",
		"owl-creek-homes-for-sale",
		"polo-fields-homes-for-sale",
	}

	for _, neighborhood := range neighborhoods {
		urls = append(urls, URL{
			Loc:        baseURL + "/" + neighborhood,
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		})
	}

	// Blog articles - get all public articles
	var articles []models.Article
	if err := db.Where("is_public = ?", true).Order("published_date DESC NULLS LAST, created_at DESC").Find(&articles).Error; err == nil {
		for _, article := range articles {
			lastMod := time.Now().Format("2006-01-02")
			if article.PublishedDate != nil {
				lastMod = article.PublishedDate.Format("2006-01-02")
			} else if !article.CreatedAt.IsZero() {
				lastMod = article.CreatedAt.Format("2006-01-02")
			}

			urls = append(urls, URL{
				Loc:        baseURL + "/blog/" + article.Slug,
				LastMod:    lastMod,
				ChangeFreq: "monthly",
				Priority:   "0.7",
			})
		}
	}

	// Create sitemap
	sitemap := Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	// Set content type to XML
	c.Header("Content-Type", "application/xml")
	c.XML(http.StatusOK, sitemap)
}
