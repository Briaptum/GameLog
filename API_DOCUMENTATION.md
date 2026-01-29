# Properties API Documentation

## Base URL
```
http://localhost:8003/api/v1
```

---

## Test Properties (Example MLS IDs)

The following MLS IDs are used in examples throughout this documentation and can be used for testing:

| MLS ID | City | Price | Beds | Baths | Status |
|--------|------|-------|------|-------|--------|
| `1691158` | Prospect, KY | $649,900 | 4 | 2.5 | Active |
| `1318142` | Taylorsville, KY | $179,900 | 3 | 2.0 | Active |
| `1690641` | Springfield, KY | $269,900 | - | - | Active |
| `1698611` | Louisville, KY | $199,900 | 4 | 2 | Active |

**To get a list of all available properties:**
```bash
# Get all properties (returns ~7500 properties)
GET /api/v1/properties/initial

# Or get a limited list with specific fields
GET /api/v1/properties?fields=listing_id,city,list_price&limit=100
```

**To test the property detail endpoint:**
```bash
# Get full details for a property
GET /api/v1/property/1691158
```

---

## Initial Properties Endpoint (Optimized for Search Page)

### **GET** `/api/v1/properties/initial`

**Purpose:** Fast endpoint to load all non-closed properties for initial page load. Returns a pre-generated, gzip-compressed cache of all active properties with essential fields.

**How it works:**
- The server pre-generates a gzip-compressed JSON file containing all non-closed properties
- This cache is refreshed automatically every hour in the background
- The endpoint serves the cached file directly (no database queries)
- If the cache is missing or corrupted, it will be regenerated on-demand before serving

**Request:**
```bash
GET /api/v1/properties/initial
```

**Response Headers:**
- `Content-Type`: `application/json`
- `Content-Encoding`: `gzip` (the response is gzip-compressed)
- `X-Cache-Generated-At`: Timestamp in RFC3339 format (e.g., `2024-12-07T15:30:45Z`) indicating when the cache was generated. This can be used by the frontend to display cache freshness or determine if a refresh is needed.

**Response Body:**
- JSON array of property objects (automatically decompressed by browser)

**Response Fields:**
Each property object contains the following fields:
- `listing_id` (string) - MLS listing ID
- `list_price` (number) - Listing price
- `city` (string) - City name
- `state_or_province` (string) - State abbreviation
- `county_or_parish` (string, nullable) - County or parish name
- `subdivision_name` (string, nullable) - Subdivision name
- `postal_code` (string, nullable) - ZIP/Postal code
- `street_number` (string) - Street number
- `street_name` (string) - Street name
- `street_suffix` (string) - Street suffix (St, Ave, etc.)
- `bedrooms_total` (number) - Number of bedrooms
- `bathrooms_total_decimal` (number) - Number of bathrooms (decimal)
- `living_area` (number) - Square footage
- `lot_size_acres` (number, nullable) - Lot size in acres
- `lot_size_square_feet` (number, nullable) - Lot size in square feet
- `basement_yn` (boolean, nullable) - Has basement (true/false)
- `property_type` (string, nullable) - Property type (e.g., "Residential", "Commercial")
- `property_sub_type` (string, nullable) - Property sub-type (e.g., "Single Family Residence")
- `latitude` (number) - Latitude coordinate
- `longitude` (number) - Longitude coordinate
- `standard_status` (string) - Property status
- `modification_timestamp` (string, nullable) - Last modification timestamp (ISO 8601 format)
- `primary_photo_url` (string, nullable) - URL to primary photo (800px width)

**Example Response:**
```json
[
  {
    "listing_id": "1691158",
    "list_price": 649900.00,
    "city": "Prospect",
    "state_or_province": "KY",
    "county_or_parish": "Jefferson",
    "subdivision_name": "Prospect Estates",
    "postal_code": "40059",
    "street_number": "123",
    "street_name": "Main",
    "street_suffix": "St",
    "bedrooms_total": 4,
    "bathrooms_total_decimal": 2.5,
    "living_area": 2500,
    "lot_size_acres": 0.5,
    "lot_size_square_feet": 21780,
    "basement_yn": true,
    "property_type": "Residential",
    "property_sub_type": "Single Family Residence",
    "latitude": 38.3456,
    "longitude": -85.6123,
    "standard_status": "Active",
    "modification_timestamp": "2024-12-05T08:24:42Z",
    "primary_photo_url": "https://cdn.resize.sparkplatform.com/lou/800x600/true/..."
  },
  {
    "listing_id": "1318142",
    "list_price": 179900.00,
    "city": "Taylorsville",
    "state_or_province": "KY",
    "county_or_parish": "Spencer",
    "subdivision_name": null,
    "postal_code": "40071",
    "street_number": "456",
    "street_name": "Oak",
    "street_suffix": "Ave",
    "bedrooms_total": 3,
    "bathrooms_total_decimal": 2.0,
    "living_area": 1800,
    "lot_size_acres": 0.25,
    "lot_size_square_feet": 10890,
    "basement_yn": false,
    "property_type": "Residential",
    "property_sub_type": "Single Family Residence",
    "latitude": 38.1234,
    "longitude": -85.4567,
    "standard_status": "Active",
    "modification_timestamp": "2024-12-04T15:30:22Z",
    "primary_photo_url": null
  }
]
```

**Frontend Usage:**
```javascript
// Fetch initial properties (browser automatically handles gzip decompression)
const response = await fetch('/api/v1/properties/initial');
const properties = await response.json();

// Get cache generation timestamp from header
const cacheGeneratedAt = response.headers.get('X-Cache-Generated-At');
if (cacheGeneratedAt) {
  const cacheDate = new Date(cacheGeneratedAt);
  const now = new Date();
  const ageInHours = (now - cacheDate) / (1000 * 60 * 60);
  
  console.log('Cache generated at:', cacheDate.toLocaleString());
  console.log('Cache age:', ageInHours.toFixed(1), 'hours');
  
  // Display cache freshness to users
  if (ageInHours > 1) {
    console.warn('Cache is older than 1 hour, consider refreshing');
  }
  
  // Use this to determine if refresh is needed
  // Cache refreshes hourly, so if older than 2 hours, might want to refresh
  if (ageInHours > 2) {
    // Optionally trigger a refresh or show warning to user
  }
}

// Use properties array for initial map/list display
properties.forEach(property => {
  // Display property on map or in list
  console.log(property.listing_id, property.list_price, property.city);
});
```

**Example: Checking Cache Freshness**
```javascript
async function fetchPropertiesWithCacheInfo() {
  const response = await fetch('/api/v1/properties/initial');
  const properties = await response.json();
  
  const cacheDateHeader = response.headers.get('X-Cache-Generated-At');
  const cacheDate = cacheDateHeader ? new Date(cacheDateHeader) : null;
  
  return {
    properties,
    cacheGeneratedAt: cacheDate,
    isFresh: cacheDate ? (Date.now() - cacheDate.getTime()) < 3600000 : false // Fresh if < 1 hour old
  };
}
```

**Important Notes:**
- This endpoint returns **all non-closed, non-rental properties** (typically ~7500 properties)
- The response is gzip-compressed, so the actual payload size is ~200-500KB
- Properties with `standard_status = 'Closed'` are excluded
- Properties with `property_type = 'Residential Lease'` (rentals) are excluded
- The cache is refreshed hourly, so data may be up to 1 hour old
- The `X-Cache-Generated-At` header indicates when the cache was last generated (useful for displaying cache freshness)
- If the cache fails to generate, the endpoint will attempt to regenerate it synchronously (may cause slight delay)
- This endpoint is optimized for initial page load - use `/api/v1/properties` for filtered searches
- Properties are ordered by `listing_id` in the cache, but can be sorted by `modification_timestamp` on the frontend

---

## Properties Query Endpoint

### **GET** or **POST** `/api/v1/properties`

Query properties with dynamic field selection, filtering, and optional media (photos, videos, etc.) inclusion.

---

## Request Parameters

### Query Parameters (GET)

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `fields` | string (comma-separated) | Property fields to return | `id,listing_id,list_price,city` |
| `include_photos` | boolean | Include media (photos, videos, etc.) with properties | `true` or `false` |
| `photo_fields` | string (comma-separated) | Media fields to return (if include_photos=true) | `media_key,uri_thumb,uri_800` |
| `limit` | integer | Maximum number of results | `50` |
| `offset` | integer | Pagination offset | `0` |
| `order_by` | string | SQL ORDER BY clause | `list_price DESC` |
| `north` | float | North latitude boundary (for bounds) | `40.7580` |
| `south` | float | South latitude boundary (for bounds) | `40.7128` |
| `east` | float | East longitude boundary (for bounds) | `-73.9352` |
| `west` | float | West longitude boundary (for bounds) | `-74.0059` |
| `{field_name}` | string | Simple equality filter | `listing_id=xyz123` |

### JSON Body (POST)

```json
{
  "fields": ["id", "listing_id", "list_price"],
  "include_photos": true,
  "photo_fields": ["media_key", "uri_thumb"],
  "filters": {
    "field_name": {
      "eq": "value",
      "gt": 100,
      "lt": 500
    }
  },
  "bounds": {
    "north": 40.7580,
    "south": 40.7128,
    "east": -73.9352,
    "west": -74.0059
  },
  "limit": 50,
  "offset": 0,
  "order_by": "list_price DESC"
}
```

---

## Filter Operators

| Operator | Description | Type | Example |
|----------|-------------|------|---------|
| `eq` | Equals | any | `{"eq": "Seattle"}` |
| `ne` | Not equals | any | `{"ne": "Portland"}` |
| `gt` | Greater than | number/date | `{"gt": 100000}` |
| `gte` | Greater than or equal | number/date | `{"gte": 100000}` |
| `lt` | Less than | number/date | `{"lt": 500000}` |
| `lte` | Less than or equal | number/date | `{"lte": 500000}` |
| `like` | Pattern match (case-sensitive) | string | `{"like": "%street%"}` |
| `ilike` | Pattern match (case-insensitive) | string | `{"ilike": "%SEATTLE%"}` |
| `in` | In array | array | `{"in": ["WA", "OR", "CA"]}` |
| `not_in` | Not in array | array | `{"not_in": ["NY", "FL"]}` |
| `is_null` | Is NULL | boolean | `{"is_null": true}` |
| `is_not_null` | Is NOT NULL | boolean | `{"is_not_null": true}` |

**Note:** Multiple operators on the same field are combined with AND logic.

---

## Bounding Box Search

The `bounds` parameter allows you to search for properties within a geographic rectangle, perfect for map viewport searches. When using bounds:

- All four values (`north`, `south`, `east`, `west`) are required
- Properties with NULL latitude or longitude are automatically excluded
- Bounds can be combined with other filters (all conditions use AND logic)
- Bounds are validated: `south < north` and `west < east`

**Bounds Parameters:**
- `north` (float): Northern latitude boundary (top of the box)
- `south` (float): Southern latitude boundary (bottom of the box)
- `east` (float): Eastern longitude boundary (right side of the box)
- `west` (float): Western longitude boundary (left side of the box)

**Example from Google Maps:**
```javascript
const bounds = map.getBounds();
const boundsData = {
  north: bounds.getNorthEast().lat(),
  south: bounds.getSouthWest().lat(),
  east: bounds.getNorthEast().lng(),
  west: bounds.getSouthWest().lng()
};
```

---

## Response Format

```json
{
  "data": [
    {
      "field1": "value1",
      "field2": "value2",
      "media": [...]  // Only if include_photos=true (returns media array, not photos)
    }
  ],
  "count": 10,
  "limit": 50,
  "offset": 0,
  "total_count": 150
}
```

---

## Examples

### Example 1: Basic Query - Select Specific Fields

**Request:**
```bash
GET /api/v1/properties?fields=id,listing_id,list_price,city,state_or_province&limit=5
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "list_price": "NjQ5OTAwLjAw",
      "city": "Prospect",
      "state_or_province": "KY"
    },
    {
      "id": 25757,
      "listing_id": "1318142",
      "list_price": "MTc5OTAwLjAw",
      "city": "Taylorsville",
      "state_or_province": "KY"
    }
  ],
  "count": 5,
  "limit": 5,
  "total_count": 5
}
```

---

### Example 2: Filter by Listing ID

**Request:**
```bash
GET /api/v1/properties?fields=id,listing_id,city,list_price&listing_id=1691158
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "city": "Prospect",
      "list_price": "NjQ5OTAwLjAw"
    }
  ],
  "count": 1,
  "total_count": 1
}
```

---

### Example 3: Price Range Filter (JSON Body)

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "list_price", "city", "state_or_province"],
  "filters": {
    "list_price": {
      "gte": 200000,
      "lte": 400000
    }
  },
  "limit": 10,
  "order_by": "list_price ASC"
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 182,
      "listing_id": "1690641",
      "list_price": "MjY5OTAwLjAw",
      "city": "Springfield",
      "state_or_province": "KY"
    }
  ],
  "count": 10,
  "limit": 10,
  "total_count": 25
}
```

---

### Example 4: Multiple Filters with City and State

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "city", "state_or_province", "list_price", "bedrooms_total", "bathrooms_full"],
  "filters": {
    "city": {
      "in": ["Louisville", "Lexington", "Frankfort"]
    },
    "state_or_province": {
      "eq": "KY"
    },
    "bedrooms_total": {
      "gte": 3
    }
  },
  "limit": 20,
  "order_by": "list_price DESC"
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 265,
      "listing_id": "1698611",
      "city": "Louisville",
      "state_or_province": "KY",
      "list_price": "MTk5OTAwLjAw",
      "bedrooms_total": 4,
      "bathrooms_full": 2
    }
  ],
  "count": 20,
  "limit": 20,
  "total_count": 45
}
```

---

### Example 5: Date Range Filter

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "created_at", "updated_at", "city"],
  "filters": {
    "created_at": {
      "gte": "2024-01-01T00:00:00Z",
      "lte": "2024-12-31T23:59:59Z"
    }
  },
  "limit": 50,
  "order_by": "created_at DESC"
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "created_at": "2024-06-30T19:09:00Z",
      "updated_at": "2024-09-25T02:38:28Z",
      "city": "Prospect"
    }
  ],
  "count": 50,
  "limit": 50,
  "total_count": 120
}
```

---

### Example 6: Geographic Bounding Box (Map Viewport Search)

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "city", "latitude", "longitude", "list_price"],
  "bounds": {
    "north": 38.5,
    "south": 38.0,
    "east": -85.0,
    "west": -85.5
  },
  "limit": 100
}
```

**Alternative (Query Parameters):**
```bash
GET /api/v1/properties?fields=id,listing_id,city,latitude,longitude,list_price&north=38.5&south=38.0&east=-85.0&west=-85.5&limit=100
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "city": "Prospect",
      "latitude": 38.3456,
      "longitude": -85.6123,
      "list_price": "NjQ5OTAwLjAw"
    }
  ],
  "count": 100,
  "limit": 100,
  "total_count": 250
}
```

**Note:** The `bounds` parameter automatically excludes properties with NULL latitude or longitude. This is ideal for map-based searches where you want to show only properties within the visible viewport.

---

### Example 7: Include Media (All Fields)

**Request:**
```bash
GET /api/v1/properties?fields=id,listing_id,city&include_photos=true&limit=2
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "city": "Prospect",
      "media": [
        {
          "id": 2182387,
          "property_id": 1446,
          "listing_id": "1691158",
          "media_key": "20250630190858317851000000",
          "media_category": "Photo",
          "name": "Primary Photo",
          "short_description": null,
          "long_description": null,
          "privacy": "Public",
          "current_privacy": "Public",
          "resource_uri": "/v1/listings/.../media/...",
          "resource_record_id": null,
          "resource_record_key": null,
          "originating_system_media_key": null,
          "is_primary": true,
          "media_order": 0,
          "object_html": null,
          "uri_thumb": "https://cdn.photos.sparkplatform.com/lou/...-t.jpg",
          "uri_300": "https://cdn.photos.sparkplatform.com/lou/....jpg",
          "uri_640": "https://cdn.resize.sparkplatform.com/lou/640x480/true/...-o.jpg",
          "uri_800": "https://cdn.resize.sparkplatform.com/lou/800x600/true/...-o.jpg",
          "uri_1024": "https://cdn.resize.sparkplatform.com/lou/1024x768/true/...-o.jpg",
          "uri_1280": "https://cdn.resize.sparkplatform.com/lou/1280x1024/true/...-o.jpg",
          "uri_1600": "https://cdn.resize.sparkplatform.com/lou/1600x1200/true/...-o.jpg",
          "uri_2048": "https://cdn.resize.sparkplatform.com/lou/2048x1600/true/...-o.jpg",
          "uri_large": "https://cdn.photos.sparkplatform.com/lou/...-o.jpg",
          "modification_timestamp": null,
          "created_at": "2025-09-25T02:38:28.802Z",
          "updated_at": "2025-09-25T02:38:28.802Z"
        }
      ]
    }
  ],
  "count": 2,
  "limit": 2,
  "total_count": 2
}
```

---

### Example 8: Include Media with Selected Fields

**Request:**
```bash
GET /api/v1/properties?fields=id,listing_id,city&include_photos=true&photo_fields=media_key,uri_thumb,uri_800,is_primary&limit=1
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "city": "Prospect",
      "media": [
        {
          "listing_id": "1691158",
          "media_key": "20250630190858317851000000",
          "uri_thumb": "https://cdn.photos.sparkplatform.com/lou/...-t.jpg",
          "uri_800": "https://cdn.resize.sparkplatform.com/lou/800x600/true/...-o.jpg",
          "is_primary": true
        },
        {
          "listing_id": "1691158",
          "media_key": "20250630190859801907000000",
          "uri_thumb": "https://cdn.photos.sparkplatform.com/lou/...-t.jpg",
          "uri_800": "https://cdn.resize.sparkplatform.com/lou/800x600/true/...-o.jpg",
          "is_primary": false
        }
      ]
    }
  ],
  "count": 1,
  "limit": 1,
  "total_count": 1
}
```

---

### Example 9: Bounding Box with Additional Filters

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "city", "latitude", "longitude", "list_price", "bedrooms_total"],
  "bounds": {
    "north": 38.5,
    "south": 38.0,
    "east": -85.0,
    "west": -85.5
  },
  "filters": {
    "list_price": {
      "gte": 200000,
      "lte": 500000
    },
    "bedrooms_total": {
      "gte": 3
    }
  },
  "limit": 50,
  "order_by": "list_price ASC"
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "city": "Prospect",
      "latitude": 38.3456,
      "longitude": -85.6123,
      "list_price": "MjY5OTAwLjAw",
      "bedrooms_total": 4
    }
  ],
  "count": 50,
  "limit": 50,
  "total_count": 87
}
```

**Note:** Bounds can be combined with other filters. All conditions are combined with AND logic.

---

### Example 10: Complex Query with Media (JSON Body)

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "city", "state_or_province", "list_price", "bedrooms_total", "bathrooms_full", "living_area"],
  "include_photos": true,
  "photo_fields": ["media_key", "uri_thumb", "uri_800", "is_primary", "media_order"],
  "filters": {
    "city": {
      "ilike": "%louisville%"
    },
    "list_price": {
      "gte": 150000,
      "lte": 300000
    },
    "bedrooms_total": {
      "gte": 3
    },
    "living_area": {
      "is_not_null": true
    }
  },
  "limit": 25,
  "offset": 0,
  "order_by": "list_price ASC"
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 265,
      "listing_id": "1698611",
      "city": "Louisville",
      "state_or_province": "KY",
      "list_price": "MTk5OTAwLjAw",
      "bedrooms_total": 4,
      "bathrooms_full": 2,
      "living_area": 2500,
      "media": [
        {
          "listing_id": "1698611",
          "media_key": "20250630190858317851000000",
          "uri_thumb": "https://cdn.photos.sparkplatform.com/lou/...-t.jpg",
          "uri_800": "https://cdn.resize.sparkplatform.com/lou/800x600/true/...-o.jpg",
          "is_primary": true,
          "media_order": 0
        }
      ]
    }
  ],
  "count": 25,
  "limit": 25,
  "offset": 0,
  "total_count": 87
}
```

---

### Example 11: String Pattern Matching

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "city", "street_name"],
  "filters": {
    "street_name": {
      "ilike": "%main%"
    },
    "city": {
      "in": ["Louisville", "Lexington"]
    }
  },
  "limit": 10
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 123,
      "listing_id": "123456",
      "city": "Louisville",
      "street_name": "Main Street"
    }
  ],
  "count": 10,
  "limit": 10,
  "total_count": 15
}
```

---

### Example 12: Null Checks

**Request:**
```bash
POST /api/v1/properties
Content-Type: application/json

{
  "fields": ["id", "listing_id", "latitude", "longitude"],
  "filters": {
    "latitude": {
      "is_not_null": true
    },
    "longitude": {
      "is_not_null": true
    }
  },
  "limit": 50
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 1446,
      "listing_id": "1691158",
      "latitude": 38.3456,
      "longitude": -85.6123
    }
  ],
  "count": 50,
  "limit": 50,
  "total_count": 1250
}
```

---

### Example 13: Pagination

**Request:**
```bash
GET /api/v1/properties?fields=id,listing_id,city&limit=10&offset=20
```

**Response:**
```json
{
  "data": [
    {
      "id": 21,
      "listing_id": "1234567",
      "city": "Louisville"
    }
  ],
  "count": 10,
  "limit": 10,
  "offset": 20,
  "total_count": 150
}
```

---

## Rental Properties Search Endpoint

### **POST** `/api/v1/properties/rentals/search`

**Purpose:** Search for rental properties by a list of unparsed addresses. This endpoint groups properties by `unparsed_address` and returns only the **latest modified** property for each address. This is optimized for rental listings pages where you need to fetch the most recent version of each property.

**Key Features:**
- Groups by `unparsed_address` using `DISTINCT ON`
- Returns only the latest version per address (ordered by `modification_timestamp` DESC)
- Perfect for rental pages where multiple listings might exist for the same address

**Request:**
```bash
POST /api/v1/properties/rentals/search
Content-Type: application/json
```

**Request Body:**
```json
{
  "addresses": [
    "4214 Winchester Rd, Louisville, KY 40207",
    "3834 Ormond Rd, Louisville, KY 40207"
  ],
  "fields": ["listing_id", "unparsed_address", "list_price", "city", "state_or_province"],
  "include_photos": true,
  "photo_fields": ["uri_thumb", "uri_800", "is_primary"]
}
```

**Request Parameters:**

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `addresses` | array of strings | Yes | List of unparsed addresses to search for | `["123 Main St, City, State 12345"]` |
| `fields` | array of strings | No | Property fields to return. If omitted, returns all fields. | `["listing_id", "list_price", "city"]` |
| `include_photos` | boolean | No | Include media (photos, videos, etc.) with properties | `true` or `false` |
| `photo_fields` | array of strings | No | Media fields to return (if include_photos=true). If omitted, returns all media fields. | `["uri_thumb", "uri_800", "is_primary"]` |

**Response Format:**
```json
{
  "data": [
    {
      "listing_id": "1382794",
      "unparsed_address": "4214 Winchester Rd, Louisville, KY 40207",
      "list_price": 1800.00,
      "city": "Louisville",
      "state_or_province": "KY",
      "bedrooms_total": 3,
      "bathrooms_total_decimal": 3.00,
      "living_area": 2157.00,
      "modification_timestamp": "2025-01-28T14:49:19Z",
      "media": [
        {
          "listing_id": "1382794",
          "uri_thumb": "https://cdn.photos.sparkplatform.com/lou/...-t.jpg",
          "uri_800": "https://cdn.resize.sparkplatform.com/lou/800x600/true/...-o.jpg",
          "is_primary": true
        }
      ]
    }
  ],
  "count": 1,
  "total_count": 1
}
```

**Response Fields:**
- `data` (array) - Array of property objects, one per unique `unparsed_address` (latest version only)
- `count` (integer) - Number of properties returned in this response
- `total_count` (integer) - Total number of unique addresses found (same as count)

**Example Request:**
```bash
POST /api/v1/properties/rentals/search
Content-Type: application/json

{
  "addresses": [
    "4214 Winchester Rd, Louisville, KY 40207",
    "3834 Ormond Rd, Louisville, KY 40207"
  ],
  "fields": [
    "listing_id",
    "unparsed_address",
    "list_price",
    "city",
    "state_or_province",
    "bedrooms_total",
    "bathrooms_total_decimal",
    "living_area",
    "modification_timestamp"
  ],
  "include_photos": true,
  "photo_fields": ["uri_thumb", "uri_800", "is_primary"]
}
```

**Example Response:**
```json
{
  "data": [
    {
      "listing_id": "1380169",
      "unparsed_address": "3834 Ormond Rd, Louisville, KY 40207",
      "list_price": 1300.00,
      "city": "Louisville",
      "state_or_province": "KY",
      "bedrooms_total": 4,
      "bathrooms_total_decimal": 2.00,
      "living_area": 1600.00,
      "modification_timestamp": "2025-01-28T14:49:16Z",
      "media": [
        {
          "listing_id": "1380169",
          "uri_thumb": "https://cdn.photos.sparkplatform.com/lou/...-t.jpg",
          "uri_800": "https://cdn.resize.sparkplatform.com/lou/800x600/true/...-o.jpg",
          "is_primary": true
        }
      ]
    },
    {
      "listing_id": "1382794",
      "unparsed_address": "4214 Winchester Rd, Louisville, KY 40207",
      "list_price": 1800.00,
      "city": "Louisville",
      "state_or_province": "KY",
      "bedrooms_total": 3,
      "bathrooms_total_decimal": 3.00,
      "living_area": 2157.00,
      "modification_timestamp": "2025-01-28T14:49:19Z",
      "media": [
        {
          "listing_id": "1382794",
          "uri_thumb": "https://cdn.photos.sparkplatform.com/lou/...-t.jpg",
          "uri_800": "https://cdn.resize.sparkplatform.com/lou/800x600/true/...-o.jpg",
          "is_primary": true
        }
      ]
    }
  ],
  "count": 2,
  "total_count": 2
}
```

**Frontend Usage:**
```javascript
// Fetch latest rental properties by addresses
async function fetchRentalProperties(addresses) {
  const response = await fetch('/api/v1/properties/rentals/search', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      addresses: addresses,
      fields: [
        'listing_id',
        'unparsed_address',
        'list_price',
        'city',
        'state_or_province',
        'bedrooms_total',
        'bathrooms_total_decimal',
        'living_area',
        'modification_timestamp',
        'latitude',
        'longitude'
      ],
      include_photos: true,
      photo_fields: ['uri_thumb', 'uri_800', 'is_primary']
    })
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = await response.json();
  return data.data; // Array of properties (latest version per address)
}

// Example usage
const rentalAddresses = [
  '4214 Winchester Rd, Louisville, KY 40207',
  '3834 Ormond Rd, Louisville, KY 40207'
];

const properties = await fetchRentalProperties(rentalAddresses);
properties.forEach(property => {
  console.log(`${property.unparsed_address}: $${property.list_price} (Latest: ${property.modification_timestamp})`);
});
```

**Important Notes:**
- **Groups by `unparsed_address`**: If multiple properties exist for the same address, only the **latest modified** one is returned
- **Latest version selection**: Properties are ordered by `modification_timestamp` DESC (or `created_at` DESC if timestamp is NULL)
- Address matching is case-sensitive and exact match only
- If an address doesn't match any property, it will not appear in the results (no error)
- The endpoint returns one property per unique `unparsed_address` in the request
- If no addresses are provided, the endpoint returns a 400 Bad Request error
- All addresses must be provided in the request body (not as query parameters)
- The endpoint supports all the same field selection and media inclusion options as the main properties endpoint
- This endpoint is specifically optimized for rental listings pages where you need the latest version of each property

**Error Responses:**

**400 Bad Request** - Missing addresses:
```json
{
  "error": "At least one address must be provided"
}
```

**400 Bad Request** - Invalid request body:
```json
{
  "error": "Invalid request body: <error message>"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Failed to query properties: <error message>"
}
```

---

## Available Property Fields

The properties table has **176 columns**. All fields listed below can be used in the `fields` parameter for selecting specific data, or in `filters` for querying. Fields are organized by category:

### Identifiers & System Fields
- `id` - Property ID (integer)
- `listing_id` - MLS Listing ID (string) - **Primary identifier**
- `listing_key` - Listing key (string, nullable)
- `originating_system_id` - Originating system ID (string, nullable)
- `originating_system_key` - Originating system key (string, nullable)
- `originating_system_name` - Originating system name (string, nullable)
- `source_system_id` - Source system ID (string, nullable)
- `source_system_key` - Source system key (string, nullable)
- `source_system_name` - Source system name (string, nullable)
- `created_at` - Creation timestamp (datetime)
- `updated_at` - Last update timestamp (datetime)

### Location & Address
- `city` - City name (string, nullable)
- `state_or_province` - State/Province abbreviation (string, nullable)
- `postal_code` - ZIP/Postal code (string, nullable)
- `county_or_parish` - County or parish name (string, nullable)
- `street_number` - Street number (string, nullable)
- `street_number_numeric` - Street number as integer (integer, nullable)
- `street_dir_prefix` - Street direction prefix (N, S, E, W) (string, nullable)
- `street_name` - Street name (string, nullable)
- `street_suffix` - Street suffix (St, Ave, Blvd, etc.) (string, nullable)
- `unit_number` - Unit number (string, nullable)
- `unparsed_address` - Full unparsed address (string, nullable)
- `subdivision_name` - Subdivision name (string, nullable)
- `high_school_district` - High school district (string, nullable)
- `latitude` - Latitude coordinate (float, nullable)
- `longitude` - Longitude coordinate (float, nullable)
- `directions` - Directions to property (string, nullable)

### Pricing
- `list_price` - Current list price (float, nullable)
- `original_list_price` - Original list price (float, nullable)
- `close_price` - Close/sale price (float, nullable)
- `price_change_timestamp` - Price change timestamp (string, nullable)

### Property Type & Status
- `property_type` - Property type (e.g., "Residential", "Commercial") (string, nullable)
- `property_sub_type` - Property sub-type (e.g., "Single Family Residence") (string, nullable)
- `standard_status` - Standard status (e.g., "Active", "Closed", "Pending") (string, nullable)
- `mls_status` - MLS status (string, nullable)
- `mls_area_major` - MLS area major (string, nullable)
- `mls_area_minor` - MLS area minor (string, nullable)
- `on_market_date` - Date property went on market (string, nullable)
- `off_market_date` - Date property went off market (string, nullable)
- `close_date` - Close date (string, nullable)
- `listing_contract_date` - Listing contract date (string, nullable)
- `purchase_contract_date` - Purchase contract date (string, nullable)
- `delayed_marketing_date` - Delayed marketing date (string, nullable)
- `delayed_marketing_yn` - Delayed marketing flag (boolean, nullable)
- `approval_status` - Approval status (boolean, nullable)

### Size & Dimensions
- `living_area` - Living area in square feet (float, nullable)
- `living_area_source` - Living area source (string, nullable)
- `building_area_total` - Total building area (float, nullable)
- `building_area_source` - Building area source (string, nullable)
- `above_grade_finished_area` - Above grade finished area (float, nullable)
- `below_grade_finished_area` - Below grade finished area (float, nullable)
- `lot_size_acres` - Lot size in acres (float, nullable)
- `lot_size_area` - Lot size area (float, nullable)
- `lot_size_square_feet` - Lot size in square feet (float, nullable)
- `lot_size_dimensions` - Lot size dimensions (string, nullable)
- `lot_size_units` - Lot size units (string, nullable)
- `stories` - Number of stories (float, nullable)
- `stories_total` - Total stories (string, nullable)
- `rooms_total` - Total number of rooms (float, nullable)

### Bedrooms & Bathrooms
- `bedrooms_total` - Total number of bedrooms (integer, nullable)
- `bathrooms_full` - Number of full bathrooms (integer, nullable)
- `bathrooms_half` - Number of half bathrooms (integer, nullable)
- `bathrooms_total_decimal` - Total bathrooms as decimal (e.g., 2.5) (float, nullable)
- `bathrooms_total_integer` - Total bathrooms as integer (float, nullable)

### Year Built & Age
- `year_built` - Year built (float, nullable)

### Features & Amenities (Boolean Flags)
- `basement_yn` - Has basement (boolean, nullable)
- `cooling_yn` - Has cooling (boolean, nullable)
- `fireplace_yn` - Has fireplace (boolean, nullable)
- `garage_yn` - Has garage (boolean, nullable)
- `heating_yn` - Has heating (boolean, nullable)
- `association_yn` - Has association (boolean, nullable)
- `internet_address_display_yn` - Internet address display flag (boolean, nullable)
- `internet_automated_valuation_display_yn` - Internet automated valuation display (boolean, nullable)
- `internet_consumer_comment_yn` - Internet consumer comment flag (boolean, nullable)
- `internet_entire_listing_display_yn` - Internet entire listing display (boolean, nullable)

### Features & Amenities (Descriptive)
- `basement` - Basement description (string, nullable)
- `cooling` - Cooling description (string, nullable)
- `heating` - Heating description (string, nullable)
- `roof` - Roof description (string, nullable)
- `architectural_style` - Architectural style (string, nullable)
- `construction_materials` - Construction materials (string, nullable)
- `exterior_features` - Exterior features (string, nullable)
- `interior_features` - Interior features (string, nullable)
- `lot_features` - Lot features (string, nullable)
- `parking_features` - Parking features (string, nullable)
- `patio_and_porch_features` - Patio and porch features (string, nullable)
- `pool_features` - Pool features (string, nullable)
- `waterfront_features` - Waterfront features (string, nullable)
- `other_structures` - Other structures (string, nullable)
- `fencing` - Fencing description (string, nullable)
- `laundry_features` - Laundry features (string, nullable)
- `levels` - Levels description (string, nullable)
- `green_energy_generation` - Green energy generation (string, nullable)
- `utilities` - Utilities description (string, nullable)
- `water_source` - Water source (string, nullable)
- `sewer` - Sewer description (string, nullable)

### Garage & Parking
- `garage_spaces` - Number of garage spaces (float, nullable)
- `carport_spaces` - Carport spaces (string, nullable)

### Fireplaces
- `fireplaces_total` - Total number of fireplaces (float, nullable)
- `fireplaces_co_basement` - Fireplaces in basement (integer, nullable)
- `fireplaces_co_level_sp_1` - Fireplaces on level 1 (integer, nullable)
- `fireplaces_co_level_sp_2` - Fireplaces on level 2 (integer, nullable)
- `fireplaces_co_level_sp_3` - Fireplaces on level 3 (integer, nullable)

### Closets
- `closets_co_basement_2` - Closets in basement (integer, nullable)
- `closets_co_level_sp_12` - Closets on level 1/2 (integer, nullable)
- `closets_co_level_sp_22` - Closets on level 2/2 (integer, nullable)
- `closets_co_level_sp_32` - Closets on level 3/2 (integer, nullable)

### Rooms (Level Information)
- `rooms_co_bedroom_sp_level` - Bedroom level (string, nullable)
- `rooms_co_dining_sp_room_sp_level` - Dining room level (string, nullable)
- `rooms_co_full_sp_bathroom_sp_level` - Full bathroom level (string, nullable)
- `rooms_co_great_sp_room_sp_level` - Great room level (string, nullable)
- `rooms_co_kitchen_sp_level` - Kitchen level (string, nullable)
- `rooms_co_laundry_sp_level_2` - Laundry level (string, nullable)
- `rooms_co_living_sp_room_sp_level` - Living room level (string, nullable)
- `rooms_co_primary_sp_bathroom_sp_level` - Primary bathroom level (string, nullable)
- `rooms_co_primary_sp_bedroom_sp_level` - Primary bedroom level (string, nullable)
- `rooms_co_rooms` - Rooms description (string, nullable)

### Property Description
- `public_remarks` - Public remarks/description (string, nullable)
- `disclosures` - Disclosures (string, nullable)
- `general_sp_property_sp_description_co_age` - Age description (string, nullable)
- `general_sp_property_sp_description_co_assumable` - Assumable description (string, nullable)
- `general_sp_property_sp_description_co_below_sp_grade_sp_unfin` - Below grade unfinished (string, nullable)
- `general_sp_property_sp_description_co_first_sp_floor_sp_pbr` - First floor description (string, nullable)
- `general_sp_property_sp_description_co_laundry_sp_level` - Laundry level description (string, nullable)
- `general_sp_property_sp_description_co_m_sp_struct_sp_flood_sp_p` - Flood plain description (string, nullable)
- `general_sp_property_sp_description_co_nonconform_sp_sq_ft_sp_fi` - Nonconforming square feet (string, nullable)
- `general_sp_property_sp_description_co_sold_sp_as_hyphen_is` - Sold as-is description (string, nullable)
- `general_sp_property_sp_description_co_sqft_sp__hyphen__sp_total` - Square feet total description (string, nullable)
- `general_sp_property_sp_description_co_total_sp__pound__sp_of_sp` - Total description (string, nullable)
- `general_sp_property_sp_description_co_total_sp_closets` - Total closets description (string, nullable)
- `general_sp_property_sp_description_co__pound__sp_1_st_sp_floor_` - 1st floor description (string, nullable)
- `general_sp_property_sp_description_co__pound__sp_2_nd_sp_floor_` - 2nd floor description (string, nullable)
- `general_sp_property_sp_description_co__pound__sp_basement_sp_be` - Basement description (string, nullable)
- `general_sp_property_sp_description_co__pound__sp_upper_sp_floor` - Upper floor description (string, nullable)

### Acreage & Land
- `acreage_sp_info_co_lake_pond_2` - Lake/pond acreage (float, nullable)
- `acreage_sp_info_co_pasture_sp_acres` - Pasture acres (float, nullable)
- `acreage_sp_info_co_tillable_sp_acres` - Tillable acres (float, nullable)
- `acreage_sp_info_co_timber_sp_acres` - Timber acres (float, nullable)
- `farm_sp_features_co_farm_sp_features` - Farm features (string, nullable)

### Association & HOA
- `association_fee` - Association/HOA fee (float, nullable)
- `association_amenities` - Association amenities (string, nullable)
- `association_fee_includes` - What association fee includes (string, nullable)
- `contract_sp_info_co_hoa_sp_fee` - HOA fee (string, nullable)

### Tax & Legal
- `parcel_number` - Parcel number (string, nullable)
- `tax_block` - Tax block (string, nullable)
- `tax_lot` - Tax lot (string, nullable)
- `location_sp_tax_sp_legal_sp_info_co_deed_sp_bk_2` - Deed book (string, nullable)
- `location_sp_tax_sp_legal_sp_info_co_disclosure_2` - Disclosure (string, nullable)
- `location_sp_tax_sp_legal_sp_info_co_pg_sp__pound__2` - Page number (string, nullable)
- `location_sp_tax_sp_legal_sp_info_co_sub_hyphen_lot_2` - Sub-lot (string, nullable)

### Media & Documents
- `photos_count` - Number of photos (float, nullable)
- `videos_count` - Number of videos (float, nullable)
- `documents_count` - Number of documents (float, nullable)
- `photos_change_timestamp` - Photos change timestamp (string, nullable)
- `floor_plans_change_timestamp` - Floor plans change timestamp (string, nullable)
- `virtual_tour_url_unbranded` - Virtual tour URL (string, nullable)

### Timestamps
- `modification_timestamp` - Last modification timestamp (string, nullable)
- `major_change_timestamp` - Major change timestamp (string, nullable)
- `major_change_type` - Major change type (string, nullable)
- `status_change_timestamp` - Status change timestamp (string, nullable)

### Listing Agent Information (Public)
- `list_agent_first_name` - Listing agent first name (string, nullable)
- `list_agent_last_name` - Listing agent last name (string, nullable)
- `list_agent_full_name` - Listing agent full name (string, nullable)
- `list_agent_mls_id` - Listing agent MLS ID (string, nullable)
- `list_agent_email` - Listing agent email (string, nullable)
- `list_agent_mobile_phone` - Listing agent mobile phone (string, nullable)
- `list_agent_aor` - Listing agent AOR (string, nullable)

### Listing Office Information
- `list_office_name` - Listing office name (string, nullable)
- `list_office_mls_id` - Listing office MLS ID (string, nullable)
- `list_office_key` - Listing office key (string, nullable)

### Co-Listing Agent Information
- `co_list_agent_first_name` - Co-listing agent first name (string, nullable)
- `co_list_agent_last_name` - Co-listing agent last name (string, nullable)
- `co_list_agent_full_name` - Co-listing agent full name (string, nullable)
- `co_list_agent_mls_id` - Co-listing agent MLS ID (string, nullable)
- `co_list_agent_key` - Co-listing agent key (string, nullable)

### Attribution
- `attribution_contact` - Attribution contact (string, nullable)

**Note:** 
- All fields are nullable (can be `null`)
- Field names are case-sensitive and use snake_case
- You can use `*` in the `fields` parameter to select all fields
- Private fields (buyer agent info, etc.) are automatically filtered out for security
- Invalid field names are silently ignored

---

## Available Media Fields

When `include_photos=true`, you can select from (response returns `media` array):

- `id` - Media ID
- `property_id` - Property ID
- `listing_id` - Listing ID
- `media_key` - Media identifier
- `media_category` - Media category (photo, video, etc.)
- `name` - Media name
- `short_description` - Short description
- `long_description` - Long description
- `privacy` - Privacy setting
- `current_privacy` - Current privacy
- `resource_uri` - Resource URI
- `resource_record_id` - Resource record ID
- `resource_record_key` - Resource record key
- `originating_system_media_key` - Originating system media key
- `is_primary` - Is primary media (boolean)
- `media_order` - Media order
- `object_html` - Object HTML (for videos/embeds)
- `uri_thumb` - Thumbnail URL
- `uri_300` - 300px URL
- `uri_640` - 640px URL
- `uri_800` - 800px URL
- `uri_1024` - 1024px URL
- `uri_1280` - 1280px URL
- `uri_1600` - 1600px URL
- `uri_2048` - 2048px URL
- `uri_large` - Large size URL
- `modification_timestamp` - Modification timestamp
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request parameters"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to query properties: <error message>"
}
```

---

## Usage Tips

1. **Initial Page Load:** Use `/api/v1/properties/initial` for fast initial property load on search pages. This endpoint returns all non-closed properties (~7500) in a gzip-compressed format (~200-500KB).
2. **Filtered Searches:** Use `/api/v1/properties` for filtered searches, pagination, and advanced queries.
3. **Rental Properties Search:** Use `/api/v1/properties/rentals/search` when you need to load the latest version of multiple rental properties by their unparsed addresses. This endpoint groups by address and returns only the most recent property per address.
4. **Performance:** Select only the fields you need using `fields` parameter
5. **Media:** Only include media when needed as it adds an additional query
6. **Pagination:** Always use `limit` to avoid large responses (not applicable to `/initial` endpoint)
7. **Filtering:** Use JSON body for complex filters with multiple operators
8. **Field Names:** Field names are case-sensitive and must match database column names (snake_case)
9. **Date Formats:** Use ISO 8601 format for dates: `"2024-01-01T00:00:00Z"`
10. **Bounding Box:** Use `bounds` for map-based searches. Properties with NULL coordinates are automatically excluded. Bounds can be combined with other filters.
11. **Cache Freshness:** The `/initial` endpoint cache is refreshed hourly. For real-time data, use `/api/v1/properties` with filters.

---

## Common Use Cases

### Load Initial Properties for Search Page
```javascript
// Fast initial load - returns all non-closed properties
const response = await fetch('/api/v1/properties/initial');
const properties = await response.json(); // Browser automatically decompresses gzip

// Display properties on map or in list
properties.forEach(property => {
  // Use property.listing_id, property.latitude, property.longitude, etc.
});
```

### Get Property by Listing ID with Media
```bash
GET /api/v1/properties?fields=id,listing_id,city,list_price&include_photos=true&photo_fields=uri_thumb,uri_800,is_primary&listing_id=1691158
```

### Search Properties in Price Range
```json
POST /api/v1/properties
{
  "fields": ["id", "listing_id", "city", "list_price", "bedrooms_total"],
  "filters": {
    "list_price": {"gte": 200000, "lte": 400000},
    "city": {"eq": "Louisville"}
  },
  "limit": 50
}
```

### Get Recent Properties with Coordinates
```json
POST /api/v1/properties
{
  "fields": ["id", "listing_id", "city", "latitude", "longitude", "created_at"],
  "filters": {
    "created_at": {"gt": "2024-01-01T00:00:00Z"},
    "latitude": {"is_not_null": true},
    "longitude": {"is_not_null": true}
  },
  "order_by": "created_at DESC",
  "limit": 100
}
```

### Search Properties in Map Viewport (Bounding Box)
```json
POST /api/v1/properties
{
  "fields": ["id", "listing_id", "city", "latitude", "longitude", "list_price"],
  "bounds": {
    "north": 40.7580,
    "south": 40.7128,
    "east": -73.9352,
    "west": -74.0059
  },
  "limit": 50
}
```

**Or via Query Parameters:**
```bash
GET /api/v1/properties?fields=id,listing_id,city,latitude,longitude,list_price&north=40.7580&south=40.7128&east=-73.9352&west=-74.0059&limit=50
```

### Load Latest Rental Properties by Addresses
```javascript
// Fetch latest rental properties by addresses (groups by address, returns latest version)
const response = await fetch('/api/v1/properties/rentals/search', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    addresses: [
      '4214 Winchester Rd, Louisville, KY 40207',
      '3834 Ormond Rd, Louisville, KY 40207'
    ],
    fields: [
      'listing_id',
      'unparsed_address',
      'list_price',
      'city',
      'state_or_province',
      'bedrooms_total',
      'bathrooms_total_decimal',
      'living_area',
      'modification_timestamp'
    ],
    include_photos: true,
    photo_fields: ['uri_thumb', 'uri_800', 'is_primary']
  })
});

const data = await response.json();
const properties = data.data; // Array of latest properties (one per address)
```

