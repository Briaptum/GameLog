package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"website-dummy/internal/services/properties_api"

	"github.com/gin-gonic/gin"
)

type PropertyController struct {
	propertiesService *properties_api.Service
}

// NewPropertyController creates a new property controller
func NewPropertyController() *PropertyController {
	return &PropertyController{
		propertiesService: properties_api.NewService(),
	}
}

// GetFeaturedProperties returns 6 random featured properties
func (pc *PropertyController) GetFeaturedProperties(c *gin.Context) {
	properties, err := pc.propertiesService.GetFeaturedProperties()
	if err != nil {
		log.Printf("Error fetching featured properties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch featured properties",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": properties,
	})
}

// QueryProperties proxies requests to the properties API
// Supports both GET (with query parameters) and POST (with JSON body)
func (pc *PropertyController) QueryProperties(c *gin.Context) {
	var queryReq properties_api.QueryRequest

	// Handle different request methods
	if c.Request.Method == "GET" {
		// Parse query parameters for GET requests
		includePhotos := c.Query("include_photos") == "true"
		
		queryReq = properties_api.QueryRequest{
			IncludePhotos: includePhotos,
		}

		// Handle listing_id filter (for single property lookup)
		if listingID := c.Query("listing_id"); listingID != "" {
			queryReq.Filters = map[string]interface{}{
				"listing_id": map[string]interface{}{
					"eq": listingID,
				},
			}
			
			// For single property lookups with photos, request all fields
			// This is for property detail pages
			// Note: Only requesting fields that we know exist in the database
			if includePhotos {
				// Don't specify fields - let the API return all available fields
				// The Properties API will return all columns it has access to
				queryReq.Fields = nil
				
				// Request all photo/media fields
				queryReq.PhotoFields = []string{
					"id", "listing_id", "media_key", "media_category", "name",
					"is_primary", "media_order",
					"uri_thumb", "uri_300", "uri_640", "uri_800", "uri_1024",
					"uri_1280", "uri_1600", "uri_2048", "uri_large",
				}
			}
		}

		// Note: For GET requests, we're keeping it simple and only supporting
		// the listing_id filter. For complex queries, use POST with JSON body.
	} else {
		// Parse JSON body for POST requests
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if err := json.Unmarshal(body, &queryReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %v", err)})
			return
		}
	}

	// Query the properties API
	response, err := pc.propertiesService.Client().QueryProperties(&queryReq)
	if err != nil {
		log.Printf("Error querying properties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to query properties: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetInitialProperties returns all non-closed properties from the cached initial endpoint
// Streams the gzip-compressed response directly from the properties API to the client
func (pc *PropertyController) GetInitialProperties(c *gin.Context) {
	// Get the streamed response from properties API (already gzip-compressed)
	resp, err := pc.propertiesService.Client().GetInitialPropertiesStream()
	if err != nil {
		log.Printf("Error fetching initial properties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch initial properties: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Copy headers from properties API response
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set status code
	c.Writer.WriteHeader(resp.StatusCode)

	// Stream the gzip-compressed body directly to the client
	// No decompression/recompression - just pass through
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		log.Printf("Error streaming properties data: %v", err)
		return
	}
}

// GetNeighborhoodProperties returns properties for a specific neighborhood
// This endpoint has hardcoded search criteria to prevent API abuse
func (pc *PropertyController) GetNeighborhoodProperties(c *gin.Context) {
	neighborhood := c.Param("neighborhood")
	
	// Build query request with hardcoded filters based on neighborhood
	var queryReq *properties_api.QueryRequest
	
	switch neighborhood {
	case "anchorage":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"city": map[string]interface{}{
					"eq": "Anchorage",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "audubon-park":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"city": map[string]interface{}{
					"eq": "Audubon Park",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "cherokee-triangle":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Cherokee Triangle",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "crescent-hill":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Crescent Hill",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "forest-springs":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"in": []string{"Forest Springs", "Forest Springs North"},
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "highlands":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Highlands",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "hunting-creek":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Hunting Creek",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "indian-hills":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Indian Hills",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "lake-forest":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Lake Forest Estates",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "mockingbird-valley":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Mockingbird Valley",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "norton-commons":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "NORTON COMMONS",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "owl-creek":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "OWL CREEK",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	case "polo-fields":
		queryReq = &properties_api.QueryRequest{
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"subdivision_name": map[string]interface{}{
					"eq": "Polo Fields",
				},
				"mls_status": map[string]interface{}{
					"in": []string{"Active", "Coming Soon", "Pending"},
				},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
		}
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Neighborhood not found",
		})
		return
	}
	
	// Query the properties API
	response, err := pc.propertiesService.Client().QueryProperties(queryReq)
	if err != nil {
		log.Printf("Error querying neighborhood properties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to query properties: %v", err),
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

