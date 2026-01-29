package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"website-dummy/internal/config"
	"website-dummy/internal/models"
	"website-dummy/internal/services/properties_api"

	"github.com/gin-gonic/gin"
)

type RentalController struct {
	propertiesService *properties_api.Service
}

// NewRentalController creates a new rental controller
func NewRentalController() *RentalController {
	return &RentalController{
		propertiesService: properties_api.NewService(),
	}
}

// GetRentals returns all rentals
// Note: Ordering by availability status requires checking properties API, which is done on frontend
// Backend returns rentals ordered by the order field
func (rc *RentalController) GetRentals(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var rentals []models.Rental

	// Use raw SQL to properly read the "order" column (reserved keyword)
	// GORM's Find() may not properly map reserved keywords even with column tags
	// Use Scan with explicit column mapping
	rows, err := db.Raw(`SELECT id, "order", unparsed_address, created_at, updated_at FROM rentals ORDER BY "order" ASC`).Rows()
	if err != nil {
		log.Printf("Error fetching rentals: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rentals",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rental models.Rental
		if err := rows.Scan(&rental.ID, &rental.Order, &rental.UnparsedAddress, &rental.CreatedAt, &rental.UpdatedAt); err != nil {
			log.Printf("Error scanning rental row: %v", err)
			continue
		}
		log.Printf("Scanned rental ID %d with order %d", rental.ID, rental.Order)
		rentals = append(rentals, rental)
	}

	// Log the orders being returned
	for _, rental := range rentals {
		log.Printf("Returning rental ID %d with order %d", rental.ID, rental.Order)
	}

	c.JSON(http.StatusOK, gin.H{
		"rentals": rentals,
	})
}

// GetPublicRentals returns all rentals for public display (no authentication required)
func (rc *RentalController) GetPublicRentals(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var rentals []models.Rental

	// Use raw SQL to properly read the "order" column (reserved keyword)
	rows, err := db.Raw(`SELECT id, "order", unparsed_address, created_at, updated_at FROM rentals ORDER BY "order" ASC`).Rows()
	if err != nil {
		log.Printf("Error fetching public rentals: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rentals",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rental models.Rental
		if err := rows.Scan(&rental.ID, &rental.Order, &rental.UnparsedAddress, &rental.CreatedAt, &rental.UpdatedAt); err != nil {
			log.Printf("Error scanning rental row: %v", err)
			continue
		}
		rentals = append(rentals, rental)
	}

	c.JSON(http.StatusOK, gin.H{
		"rentals": rentals,
	})
}

// GetPublicRentalsWithProperties returns all rentals with their property data from the properties API
// This endpoint fetches rentals from our database, then searches the properties API for each one
func (rc *RentalController) GetPublicRentalsWithProperties(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	// Fetch all rentals from database
	var rentals []models.Rental
	rows, err := db.Raw(`SELECT id, "order", unparsed_address, created_at, updated_at FROM rentals ORDER BY "order" ASC`).Rows()
	if err != nil {
		log.Printf("Error fetching public rentals: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rentals",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rental models.Rental
		if err := rows.Scan(&rental.ID, &rental.Order, &rental.UnparsedAddress, &rental.CreatedAt, &rental.UpdatedAt); err != nil {
			log.Printf("Error scanning rental row: %v", err)
			continue
		}
		rentals = append(rentals, rental)
	}

	if len(rentals) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"properties": []interface{}{},
		})
		return
	}

	// Extract addresses from rentals
	addresses := make([]string, 0, len(rentals))
	for _, rental := range rentals {
		addresses = append(addresses, rental.UnparsedAddress)
	}

	// Search properties API for these addresses
	properties, err := rc.propertiesService.Client().SearchRentals(
		addresses,
		[]string{
			"listing_id",
			"unparsed_address",
			"list_price",
			"city",
			"state_or_province",
			"bedrooms_total",
			"bathrooms_total_decimal",
			"living_area",
			"latitude",
			"longitude",
			"modification_timestamp",
			"mls_status",
		},
		true, // include photos
		[]string{"uri_thumb", "uri_800", "is_primary"},
	)

	if err != nil {
		log.Printf("Error searching rental properties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to search properties: %v", err),
		})
		return
	}

	// Include all rentals (both Active and Closed)
	// Create a map of address -> property for quick lookup
	propertyMap := make(map[string]map[string]interface{})
	for _, prop := range properties {
		if addr, ok := prop["unparsed_address"].(string); ok {
			propertyMap[addr] = prop
		// Log mls_status for debugging
		if status, ok := prop["mls_status"].(string); ok {
			log.Printf("Property %s has mls_status: %q", addr, status)
		} else if status, ok := prop["mls_status"].(interface{}); ok {
			log.Printf("Property %s has mls_status (non-string): %v (type: %T)", addr, status, status)
		} else {
			log.Printf("Property %s missing mls_status field", addr)
		}
		}
	}

	// Match properties to rentals by address and preserve order
	result := make([]map[string]interface{}, 0, len(rentals))
	for _, rental := range rentals {
		if prop, found := propertyMap[rental.UnparsedAddress]; found {
			result = append(result, prop)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"properties": result,
		"count":      len(result),
	})
}

// CreateRental handles POST requests to create a new rental
func (rc *RentalController) CreateRental(c *gin.Context) {
	var req struct {
		Order           *int   `json:"order"` // Optional, will be auto-calculated if not provided
		UnparsedAddress string `json:"unparsed_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. unparsed_address is required.",
		})
		return
	}

	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	// Auto-calculate order if not provided
	order := 0
	if req.Order != nil {
		order = *req.Order
	} else {
		// Get the maximum order value and add 1
		var maxOrder int
		if err := db.Model(&models.Rental{}).Select("COALESCE(MAX(\"order\"), 0)").Scan(&maxOrder).Error; err != nil {
			log.Printf("Error getting max order: %v", err)
			// Default to 0 if there's an error
			order = 0
		} else {
			order = maxOrder + 1
		}
	}

	rental := models.Rental{
		Order:           order,
		UnparsedAddress: req.UnparsedAddress,
	}

	if err := db.Create(&rental).Error; err != nil {
		log.Printf("Error creating rental: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create rental",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rental created successfully",
		"rental":  rental,
	})
}

// UpdateRental handles PUT requests to update a rental
func (rc *RentalController) UpdateRental(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var rental models.Rental

	// First check if the rental exists
	if err := db.First(&rental, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Rental not found",
		})
		return
	}

	var req struct {
		Order           *int    `json:"order"`
		UnparsedAddress *string `json:"unparsed_address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Update fields if provided
	if req.Order != nil {
		log.Printf("Updating rental ID %s: old order=%d, new order=%d", id, rental.Order, *req.Order)
		rental.Order = *req.Order
	}
	if req.UnparsedAddress != nil {
		rental.UnparsedAddress = *req.UnparsedAddress
	}

	// Use Update with specific column to ensure the order field is updated
	if req.Order != nil {
		if err := db.Model(&rental).Update("order", rental.Order).Error; err != nil {
			log.Printf("Error updating rental order: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update rental order",
				"details": err.Error(),
			})
			return
		}
	}
	
	if req.UnparsedAddress != nil {
		if err := db.Model(&rental).Update("unparsed_address", rental.UnparsedAddress).Error; err != nil {
			log.Printf("Error updating rental address: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update rental address",
				"details": err.Error(),
			})
			return
		}
	}

	// Reload the rental to get the updated values
	if err := db.First(&rental, id).Error; err != nil {
		log.Printf("Error reloading rental after update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reload rental after update",
		})
		return
	}

	log.Printf("✅ Updated rental ID %s with order %d", id, rental.Order)
	c.JSON(http.StatusOK, gin.H{
		"message": "Rental updated successfully",
		"rental":  rental,
	})
}

// DeleteRental handles DELETE requests to delete a rental
func (rc *RentalController) DeleteRental(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var rental models.Rental

	// First check if the rental exists
	if err := db.First(&rental, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Rental not found",
		})
		return
	}

	// Hard delete the rental
	if err := db.Delete(&rental, id).Error; err != nil {
		log.Printf("Error deleting rental ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete rental",
		})
		return
	}

	log.Printf("✅ Deleted rental ID %s", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Rental deleted successfully",
		"id":      id,
	})
}

// SearchProperties searches the properties API by address or MLS ID
func (rc *RentalController) SearchProperties(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameter 'q' is required",
		})
		return
	}

	// Determine if query is an MLS ID (numeric) or address (text)
	// Try to parse as integer to check if it's an MLS ID
	var queryReq *properties_api.QueryRequest
	
	if _, err := strconv.Atoi(query); err == nil {
		// It's numeric, treat as MLS ID
		queryReq = &properties_api.QueryRequest{
			Fields: []string{
				"id",
				"listing_id",
				"unparsed_address",
				"mls_status",
				"list_price",
				"bedrooms_total",
				"bathrooms_full",
				"city",
				"state_or_province",
			},
			Filters: map[string]interface{}{
				"listing_id": map[string]interface{}{
					"eq": query,
				},
				"property_type": map[string]interface{}{
					"eq": "Residential Lease",
				},
			},
			Limit: 50,
		}
	} else {
		// It's text, treat as address search
		queryReq = &properties_api.QueryRequest{
			Fields: []string{
				"id",
				"listing_id",
				"unparsed_address",
				"mls_status",
				"list_price",
				"bedrooms_total",
				"bathrooms_full",
				"city",
				"state_or_province",
			},
			Filters: map[string]interface{}{
				"unparsed_address": map[string]interface{}{
					"like": fmt.Sprintf("%%%s%%", strings.ReplaceAll(query, "%", "\\%")),
				},
				"property_type": map[string]interface{}{
					"eq": "Residential Lease",
				},
			},
			Limit: 50,
		}
	}

	// Query the properties API
	response, err := rc.propertiesService.Client().QueryProperties(queryReq)
	if err != nil {
		log.Printf("Error searching properties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to search properties: %v", err),
		})
		return
	}

	// Convert response data to properties using the service's conversion
	// We'll work with the raw data to extract unparsed_address
	rawProperties := response.Data

	// Format results for frontend from raw data
	results := make([]map[string]interface{}, 0)
	
	// Handle different response data types
	var rawProps []map[string]interface{}
	if maps, ok := rawProperties.([]map[string]interface{}); ok {
		rawProps = maps
	} else if interfaces, ok := rawProperties.([]interface{}); ok {
		rawProps = make([]map[string]interface{}, 0, len(interfaces))
		for _, item := range interfaces {
			if itemMap, ok := item.(map[string]interface{}); ok {
				rawProps = append(rawProps, itemMap)
			}
		}
	}
	
	for _, rawProp := range rawProps {
		// Extract fields from raw property
		listingID := ""
		if id, ok := rawProp["listing_id"].(string); ok {
			listingID = id
		}
		
		unparsedAddress := ""
		if addr, ok := rawProp["unparsed_address"].(string); ok {
			unparsedAddress = addr
		}
		
		var mlsStatusPtr *string
		if status, ok := rawProp["mls_status"].(string); ok {
			mlsStatusPtr = &status
		}
		
		// Determine if property is available
		isAvailable := false
		if mlsStatusPtr != nil {
			status := strings.ToLower(strings.TrimSpace(*mlsStatusPtr))
			isAvailable = status == "active" || 
				status == "coming soon" || 
				status == "pending" || 
				status == "active under contract"
		}
		
		result := map[string]interface{}{
			"listing_id":       listingID,
			"unparsed_address": unparsedAddress,
			"mls_status":       mlsStatusPtr,
			"is_available":     isAvailable,
		}
		
		// Add price, beds, baths if available
		if price, ok := rawProp["list_price"].(float64); ok {
			result["list_price"] = price
		}
		if beds, ok := rawProp["bedrooms_total"].(float64); ok {
			result["bedrooms_total"] = int(beds)
		}
		if baths, ok := rawProp["bathrooms_full"].(float64); ok {
			result["bathrooms_full"] = baths
		}
		
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
	})
}

// BulkUpdateRentalOrders handles PUT requests to update multiple rental orders at once
func (rc *RentalController) BulkUpdateRentalOrders(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var req struct {
		Rentals []struct {
			ID    uint `json:"id" binding:"required"`
			Order int  `json:"order" binding:"required"`
		} `json:"rentals" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format. Expected array of {id, order}",
		})
		return
	}

	// Update all rentals in a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, rentalUpdate := range req.Rentals {
		log.Printf("Updating rental ID %d order to %d", rentalUpdate.ID, rentalUpdate.Order)
		// Use raw SQL to update the reserved keyword column
		if err := tx.Exec(`UPDATE rentals SET "order" = ? WHERE id = ?`, rentalUpdate.Order, rentalUpdate.ID).Error; err != nil {
			tx.Rollback()
			log.Printf("Error updating rental ID %d order to %d: %v", rentalUpdate.ID, rentalUpdate.Order, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to update rental ID %d", rentalUpdate.ID),
				"details": err.Error(),
			})
			return
		}
		// Verify the update using raw SQL
		var orderValue int
		if err := tx.Raw(`SELECT "order" FROM rentals WHERE id = ?`, rentalUpdate.ID).Scan(&orderValue).Error; err == nil {
			log.Printf("Verified rental ID %d now has order %d", rentalUpdate.ID, orderValue)
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing bulk rental order update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit rental order updates",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Bulk updated %d rental orders", len(req.Rentals))
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully updated %d rental orders", len(req.Rentals)),
		"count":   len(req.Rentals),
	})
}


