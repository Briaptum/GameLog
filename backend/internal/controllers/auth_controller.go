package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"website-dummy/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthController handles authentication-related endpoints
type AuthController struct{}

// NewAuthController creates a new auth controller
func NewAuthController() *AuthController {
	return &AuthController{}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// User represents a user
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Sites    []string `json:"sites"`
}

// Auth0TokenResponse represents the response from Auth0 token endpoint
type Auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Auth0User represents a user from Auth0 Management API
type Auth0User struct {
	UserID      string                 `json:"user_id"`
	Email       string                 `json:"email"`
	AppMetadata map[string]interface{} `json:"app_metadata"`
}

// Auth0Error represents an error from Auth0
type Auth0Error struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Login handles custom login form submission
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get Auth0 configuration from environment
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0ClientID := os.Getenv("AUTH0_CLIENT_ID")
	auth0ClientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	if auth0Domain == "" || auth0ClientID == "" || auth0ClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth0 configuration missing"})
		return
	}

	// Step 1: Authenticate with Auth0 using Resource Owner Password Grant
	log.Printf("Attempting Auth0 authentication for user: %s", req.Email)
	log.Printf("Using domain: %s, client: %s", auth0Domain, auth0ClientID)
	_, err := ac.authenticateWithAuth0(auth0Domain, auth0ClientID, auth0ClientSecret, req.Email, req.Password)
	if err != nil {
		log.Printf("Auth0 authentication failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	log.Printf("Auth0 authentication successful")

	// Step 2: Get user information from Auth0 Management API
	log.Printf("Fetching user details from Auth0 Management API...")
	user, err := ac.getUserFromAuth0(auth0Domain, auth0ClientID, auth0ClientSecret, req.Email)
	if err != nil {
		log.Printf("Failed to get user information: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}
	log.Printf("User details fetched successfully: %s", user.Email)

	// Step 3: Check if user has required permissions
	log.Printf("Checking user permissions...")
	if !ac.hasRequiredPermissions(user) {
		log.Printf("Access denied: User %s does not have required permissions", req.Email)
		log.Printf("User app_metadata: %+v", user.AppMetadata)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Insufficient permissions"})
		return
	}
	log.Printf("User has required permissions")

	// Step 4: Create our own JWT token
	jwtToken, err := ac.createJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: jwtToken,
		User: User{
			ID:    user.UserID,
			Email: user.Email,
			Role:  ac.getUserRole(user),
			Sites: ac.getUserSites(user),
		},
	})
}

// Logout handles user logout
func (ac *AuthController) Logout(c *gin.Context) {
	auth0Domain := c.GetString("AUTH0_DOMAIN")
	clientID := c.GetString("AUTH0_CLIENT_ID")

	if auth0Domain == "" || clientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth0 configuration missing"})
		return
	}

	// Build Auth0 logout URL
	logoutURL := url.URL{
		Scheme: "https",
		Host:   auth0Domain,
		Path:   "/v2/logout",
	}

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("returnTo", "http://localhost:3000")

	logoutURL.RawQuery = params.Encode()

	c.Redirect(http.StatusTemporaryRedirect, logoutURL.String())
}

// Profile returns the current user's profile
func (ac *AuthController) Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// authenticateWithAuth0 authenticates user with Auth0 using Resource Owner Password Grant
func (ac *AuthController) authenticateWithAuth0(domain, clientID, clientSecret, email, password string) (*Auth0TokenResponse, error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)

	payload := map[string]interface{}{
		"grant_type":    "password",
		"username":      email,
		"password":      password,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"scope":         "openid profile email",
		"realm":         "Username-Password-Authentication", // Specify the database connection
	}

	// Don't include audience for password grant - it causes "default connection" errors
	// The Management API audience is only needed for M2M tokens

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var auth0Err Auth0Error
		json.NewDecoder(resp.Body).Decode(&auth0Err)
		return nil, fmt.Errorf("auth0 error: %s", auth0Err.ErrorDescription)
	}

	var tokenResp Auth0TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// getUserFromAuth0 gets user information from Auth0 Management API
func (ac *AuthController) getUserFromAuth0(domain, clientID, clientSecret, email string) (*Auth0User, error) {
	// First get a Management API token
	mgmtToken, err := ac.getManagementToken(domain, clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get management token: %w", err)
	}

	// Then get user by email using Management API
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/v2/users-by-email?email=%s", domain, url.QueryEscape(email)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user from Auth0: %s", string(body))
	}

	var users []Auth0User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

// getManagementToken gets an access token for the Management API
func (ac *AuthController) getManagementToken(domain, clientID, clientSecret string) (string, error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)

	payload := map[string]interface{}{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      fmt.Sprintf("https://%s/api/v2/", domain),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get management token: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token in response")
	}

	return accessToken, nil
}

// hasRequiredPermissions checks if user has admin role or website-dummy-com site access
func (ac *AuthController) hasRequiredPermissions(user *Auth0User) bool {
	role := ac.getUserRole(user)
	sites := ac.getUserSites(user)

	// Check for admin role
	if role == "admin" {
		return true
	}

	// Check for website-dummy-com site access
	for _, site := range sites {
		if site == "website-dummy-com" {
			return true
		}
	}

	return false
}

// getUserRole extracts role from user's app metadata
func (ac *AuthController) getUserRole(user *Auth0User) string {
	if user.AppMetadata == nil {
		log.Printf("User has no app_metadata")
		return ""
	}
	if role, ok := user.AppMetadata["role"].(string); ok {
		log.Printf("User role: %s", role)
		return role
	}
	log.Printf("No role found in app_metadata")
	return ""
}

// getUserSites extracts sites from user's app metadata
func (ac *AuthController) getUserSites(user *Auth0User) []string {
	if user.AppMetadata == nil {
		log.Printf("User has no app_metadata for sites")
		return []string{}
	}
	if sites, ok := user.AppMetadata["sites"].([]interface{}); ok {
		var siteStrings []string
		for _, site := range sites {
			if siteStr, ok := site.(string); ok {
				siteStrings = append(siteStrings, siteStr)
			}
		}
		log.Printf("User sites: %v", siteStrings)
		return siteStrings
	}
	log.Printf("No sites found in app_metadata")
	return []string{}
}

// createJWTToken creates our own JWT token with user information
func (ac *AuthController) createJWTToken(user *Auth0User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.UserID,
		"email": user.Email,
		"https://your-namespace.com/app_metadata": map[string]interface{}{
			"role":  ac.getUserRole(user),
			"sites": ac.getUserSites(user),
		},
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte("your-secret-key"))
}

// ChangePasswordRequest represents the password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword handles password change requests
func (ac *AuthController) ChangePassword(c *gin.Context) {
	// Get authenticated user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// The middleware stores user as *middleware.Auth0Claims
	claims, ok := userInterface.(*middleware.Auth0Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data format"})
		return
	}

	userID := claims.Sub
	email := claims.Email

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate new password length
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 8 characters long"})
		return
	}

	// Get Auth0 configuration
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0ClientID := os.Getenv("AUTH0_CLIENT_ID")
	auth0ClientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	if auth0Domain == "" || auth0ClientID == "" || auth0ClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth0 configuration missing"})
		return
	}

	// Step 1: Verify current password by attempting authentication
	log.Printf("Verifying current password for user: %s", userID)
	_, err := ac.authenticateWithAuth0(auth0Domain, auth0ClientID, auth0ClientSecret, email, req.CurrentPassword)
	if err != nil {
		log.Printf("Current password verification failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}
	log.Printf("Current password verified successfully")

	// Step 2: Get Management API token
	mgmtToken, err := ac.getManagementToken(auth0Domain, auth0ClientID, auth0ClientSecret)
	if err != nil {
		log.Printf("Failed to get management token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get management token"})
		return
	}

	// Step 3: Update password via Auth0 Management API
	updateURL := fmt.Sprintf("https://%s/api/v2/users/%s", auth0Domain, url.QueryEscape(userID))
	
	updatePayload := map[string]interface{}{
		"password": req.NewPassword,
	}

	jsonData, err := json.Marshal(updatePayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare update request"})
		return
	}

	httpReq, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create update request"})
		return
	}

	httpReq.Header.Set("Authorization", "Bearer "+mgmtToken)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Auth0 password update failed: %s", string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	log.Printf("Password updated successfully for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
