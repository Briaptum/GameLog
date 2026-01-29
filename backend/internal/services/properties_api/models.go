package properties_api

import (
	"encoding/json"
	"strconv"
)

// Float64OrString handles both float64 and string values from JSON
type Float64OrString float64

// UnmarshalJSON implements custom unmarshaling for Float64OrString
func (f *Float64OrString) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch val := v.(type) {
	case float64:
		*f = Float64OrString(val)
	case string:
		if val == "" {
			return nil
		}
		// Parse directly as float string
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		*f = Float64OrString(parsed)
	case nil:
		return nil
	default:
		return json.Unmarshal(data, (*float64)(f))
	}
	return nil
}

// Float64 returns the float64 value
func (f Float64OrString) Float64() float64 {
	return float64(f)
}

// IntOrString handles both int and string values from JSON
type IntOrString int

// UnmarshalJSON implements custom unmarshaling for IntOrString
func (i *IntOrString) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch val := v.(type) {
	case float64:
		*i = IntOrString(int(val))
	case int:
		*i = IntOrString(val)
	case string:
		if val == "" {
			return nil
		}
		// Try parsing as float first (might be "4.00")
		parsedFloat, err := strconv.ParseFloat(val, 64)
		if err == nil {
			*i = IntOrString(int(parsedFloat))
			return nil
		}
		// If float parse failed, try as int
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		*i = IntOrString(parsed)
	case nil:
		return nil
	default:
		return json.Unmarshal(data, (*int)(i))
	}
	return nil
}

// Int returns the int value
func (i IntOrString) Int() int {
	return int(i)
}

// Property represents a property from the API
type Property struct {
	ID              int64          `json:"id"`
	ListingID       string         `json:"listing_id"`
	StreetNumber    *string        `json:"street_number"`
	StreetName      *string        `json:"street_name"`
	StreetSuffix    *string        `json:"street_suffix"`
	City            *string        `json:"city"`
	StateOrProvince *string        `json:"state_or_province"`
	UnparsedAddress *string        `json:"unparsed_address"`
	ListPrice       *Float64OrString `json:"list_price"`
	BedsTotal       *IntOrString     `json:"bedrooms_total"`
	BathsTotal      *Float64OrString `json:"bathrooms_full"`
	LivingArea      *IntOrString      `json:"living_area"`
	Latitude        *Float64OrString `json:"latitude"`
	Longitude       *Float64OrString `json:"longitude"`
	MLSStatus       *string        `json:"mls_status"`
	Photos          []Photo        `json:"media,omitempty"` // API returns "media" array
}

// Photo represents a photo/media from the API
type Photo struct {
	PhotoID   string `json:"media_key"` // API now uses "media_key" instead of "photo_id"
	URIThumb  string `json:"uri_thumb"`
	URI640    string `json:"uri_640"`
	URI800    string `json:"uri_800"`
	URI1024   string `json:"uri_1024"`
	URI1280   string `json:"uri_1280"`
	URI1600   string `json:"uri_1600"`
	URI2048   string `json:"uri_2048"`
	URILarge  string `json:"uri_large"`
	IsPrimary bool   `json:"is_primary"`
}

// PropertiesResponse represents the API response
type PropertiesResponse struct {
	Data       interface{} `json:"data"` // Can be []Property or []map[string]interface{}
	Count      int         `json:"count"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset,omitempty"`
	TotalCount int         `json:"total_count"`
}

// Bounds represents geographic bounds for filtering properties
type Bounds struct {
	North float64 `json:"north"`
	South float64 `json:"south"`
	East  float64 `json:"east"`
	West  float64 `json:"west"`
}

// QueryRequest represents a query to the properties API
type QueryRequest struct {
	Fields        []string               `json:"fields,omitempty"`
	IncludePhotos bool                   `json:"include_photos,omitempty"`
	PhotoFields   []string               `json:"photo_fields,omitempty"`
	Filters       map[string]interface{} `json:"filters,omitempty"`
	Bounds        *Bounds                `json:"bounds,omitempty"`
	Limit         int                    `json:"limit,omitempty"`
	Offset        int                    `json:"offset,omitempty"`
	OrderBy       string                 `json:"order_by,omitempty"`
}

