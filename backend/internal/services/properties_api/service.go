package properties_api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// convertDataToProperties converts PropertiesResponse.Data (interface{}) to []Property
func convertDataToProperties(data interface{}) ([]Property, error) {
	// If it's already []Property, return it
	if props, ok := data.([]Property); ok {
		return props, nil
	}
	
	// If it's []map[string]interface{}, convert to []Property
	if maps, ok := data.([]map[string]interface{}); ok {
		props := make([]Property, 0, len(maps))
		for _, m := range maps {
			// Convert map to JSON and then to Property
			jsonData, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal property map: %w", err)
			}
			var prop Property
			if err := json.Unmarshal(jsonData, &prop); err != nil {
				return nil, fmt.Errorf("failed to unmarshal property: %w", err)
			}
			props = append(props, prop)
		}
		return props, nil
	}
	
	// If it's []interface{}, convert each element to map[string]interface{} then to Property
	if interfaces, ok := data.([]interface{}); ok {
		props := make([]Property, 0, len(interfaces))
		for _, item := range interfaces {
			// Convert interface{} to map[string]interface{}
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected element type in []interface{}: %T", item)
			}
			
			// Convert map to JSON and then to Property
			jsonData, err := json.Marshal(itemMap)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal property map: %w", err)
			}
			var prop Property
			if err := json.Unmarshal(jsonData, &prop); err != nil {
				return nil, fmt.Errorf("failed to unmarshal property: %w", err)
			}
			props = append(props, prop)
		}
		return props, nil
	}
	
	return nil, fmt.Errorf("unexpected data type: %T", data)
}

// Service provides high-level methods for working with properties
type Service struct {
	client *Client
}

// NewService creates a new properties service
func NewService() *Service {
	return &Service{
		client: NewClient(),
	}
}

// Client returns the underlying client for direct API access
func (s *Service) Client() *Client {
	return s.client
}

// GetFeaturedProperties returns the top 6 most expensive properties for sale
func (s *Service) GetFeaturedProperties() ([]FeaturedProperty, error) {
	// First, get properties for agent 16619
	agentReq := &QueryRequest{
		Fields: []string{
			"id",
			"listing_id",
			"street_number",
			"street_name",
			"street_suffix",
			"city",
			"state_or_province",
			"list_price",
			"bedrooms_total",
			"bathrooms_full",
			"living_area",
			"mls_status",
		},
		IncludePhotos: true,
		PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
		Filters: map[string]interface{}{
			"list_office_mls_id": map[string]interface{}{
				"eq": "16619",
			},
			"mls_status": map[string]interface{}{
				"in": []string{"Active", "Coming Soon", "Pending", "Active Under Contract"},
			},
			"property_type": map[string]interface{}{
				"in": []string{"Residential"},
			},
		},
		Limit:   6,                    // Get up to 6 most expensive
		OrderBy: "list_price DESC", // Order by price descending
	}

	agentResp, err := s.client.QueryProperties(agentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to query agent properties: %w", err)
	}

	properties, err := convertDataToProperties(agentResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert properties data: %w", err)
	}
	propertyIDs := make(map[int64]bool)
	for _, prop := range properties {
		propertyIDs[prop.ID] = true
	}

	// If we don't have 6 properties, fill with Anchorage properties
	if len(properties) < 6 {
		needed := 6 - len(properties)
		
		anchorageReq := &QueryRequest{
			Fields: []string{
				"id",
				"listing_id",
				"street_number",
				"street_name",
				"city",
				"state_or_province",
				"list_price",
				"bedrooms_total",
				"bathrooms_full",
				"living_area",
				"mls_status",
			},
			IncludePhotos: true,
			PhotoFields:   []string{"media_key", "uri_640", "is_primary"},
			Filters: map[string]interface{}{
				"city": map[string]interface{}{
					"eq": "Anchorage",
				},
			"mls_status": map[string]interface{}{
				"in": []string{"Active", "Coming Soon", "Pending", "Active Under Contract"},
			},
				"property_type": map[string]interface{}{
					"in": []string{"Residential"},
				},
			},
			Limit:   needed,
			OrderBy: "list_price DESC",
		}

		anchorageResp, err := s.client.QueryProperties(anchorageReq)
		if err == nil {
			anchorageProperties, convErr := convertDataToProperties(anchorageResp.Data)
			if convErr == nil && len(anchorageProperties) > 0 {
				// Add Anchorage properties that aren't already in our list
				for _, prop := range anchorageProperties {
					if !propertyIDs[prop.ID] && len(properties) < 6 {
						properties = append(properties, prop)
						propertyIDs[prop.ID] = true
					}
				}
			}
		}
	}

	if len(properties) == 0 {
		return []FeaturedProperty{}, nil
	}

	// Convert to featured properties format
	featured := make([]FeaturedProperty, 0, len(properties))
	for _, prop := range properties {
		featuredProp := s.convertToFeatured(prop)
		if featuredProp != nil {
			featured = append(featured, *featuredProp)
		}
	}

	return featured, nil
}

// FeaturedProperty represents a property formatted for the featured section
type FeaturedProperty struct {
	ListingID string   `json:"listing_id"`
	Address   string   `json:"address"`
	Location  string   `json:"location"`
	Price     string   `json:"price"`
	Image     string   `json:"image"`
	Beds      *int     `json:"beds,omitempty"`
	Baths     *float64 `json:"baths,omitempty"`
	Sqft      *int     `json:"sqft,omitempty"`
	Status    *string  `json:"status,omitempty"` // "Pending", "Coming Soon", or "Active Under Contract"
}

// convertToFeatured converts an API property to a featured property
func (s *Service) convertToFeatured(prop Property) *FeaturedProperty {
	// Build address
	addressParts := []string{}
	if prop.StreetNumber != nil {
		addressParts = append(addressParts, *prop.StreetNumber)
	}
	if prop.StreetName != nil {
		addressParts = append(addressParts, *prop.StreetName)
	}
	if prop.StreetSuffix != nil {
		addressParts = append(addressParts, *prop.StreetSuffix)
	}
	address := strings.Join(addressParts, " ")
	if address == "" {
		address = "Address not available"
	}

	// Build location
	locationParts := []string{}
	if prop.City != nil {
		locationParts = append(locationParts, *prop.City)
	}
	if prop.StateOrProvince != nil {
		locationParts = append(locationParts, *prop.StateOrProvince)
	}
	location := strings.Join(locationParts, ", ")
	if location == "" {
		location = "Location not available"
	}

	// Format price
	price := "Price not available"
	if prop.ListPrice != nil {
		priceFloat := prop.ListPrice.Float64()
		price = formatPrice(priceFloat)
	}

	// Get primary photo or first photo, using 640px for featured properties
	image := ""
	getPhotoURI := func(photo Photo) string {
		// Use 640px (max size for featured properties)
		if photo.URI640 != "" {
			return photo.URI640
		}
		// Fallback to 800px if 640px not available
		if photo.URI800 != "" {
			return photo.URI800
		}
		return ""
	}
	
	for _, photo := range prop.Photos {
		if photo.IsPrimary {
			image = getPhotoURI(photo)
			if image != "" {
				break
			}
		}
	}
	if image == "" && len(prop.Photos) > 0 {
		image = getPhotoURI(prop.Photos[0])
	}

	// Determine status badge - API returns exactly "Pending", "Coming Soon", or "Active Under Contract"
	var status *string
	if prop.MLSStatus != nil {
		statusLower := strings.ToLower(strings.TrimSpace(*prop.MLSStatus))
		
		if statusLower == "pending" {
			statusVal := "Pending"
			status = &statusVal
		} else if statusLower == "coming soon" {
			statusVal := "Coming Soon"
			status = &statusVal
		} else if statusLower == "active under contract" {
			statusVal := "Active Under Contract"
			status = &statusVal
		}
	}

	// Convert beds_total
	var beds *int
	if prop.BedsTotal != nil {
		bedsVal := prop.BedsTotal.Int()
		beds = &bedsVal
	}

	// Convert baths_total
	var baths *float64
	if prop.BathsTotal != nil {
		bathsVal := prop.BathsTotal.Float64()
		baths = &bathsVal
	}

	// Convert living_area
	var sqft *int
	if prop.LivingArea != nil {
		sqftVal := prop.LivingArea.Int()
		sqft = &sqftVal
	}

	return &FeaturedProperty{
		ListingID: prop.ListingID,
		Address:   address,
		Location:  location,
		Price:     price,
		Image:     image,
		Beds:      beds,
		Baths:     baths,
		Sqft:      sqft,
		Status:    status,
	}
}

// formatPrice formats a price as currency with commas
func formatPrice(price float64) string {
	priceInt := int64(price)
	
	// Format with commas
	priceStr := fmt.Sprintf("%d", priceInt)
	n := len(priceStr)
	if n <= 3 {
		return fmt.Sprintf("$%s", priceStr)
	}
	
	// Add commas every 3 digits from right
	result := ""
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(priceStr[i])
	}
	
	return fmt.Sprintf("$%s", result)
}

