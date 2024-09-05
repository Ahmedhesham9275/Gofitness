package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myblog/config"
	"myblog/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to register a user
func registerUser(username, password string, router http.Handler) *httptest.ResponseRecorder {

	jsonValue, _ := json.Marshal(map[string]string{
		"Username": username,
		"Password": password,
	})
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, req)

	return writer
}

// Helper function to login and retrieve the token
func loginUser(username, password string, router http.Handler) (string, *httptest.ResponseRecorder) {

	jsonValue, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, req)

	// Parse the token from the response body
	var responseData map[string]string
	json.NewDecoder(writer.Body).Decode(&responseData)

	token := responseData["token"]

	return token, writer
}

func TestRegister(t *testing.T) {
	clearTables()

	router := setupRouter()

	// Use the helper function to register the user
	w := registerUser("username", "password", router)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	clearTables()

	router := setupRouter()

	// Register the user using the helper function
	w := registerUser("testuser", "password", router)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Login using the helper function
	token, w := loginUser("testuser", "password", router)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, token)

	// Login with the wrong password to assert login fail
	_, w = loginUser("testuser", "wrongpassword", router)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPost(t *testing.T) {
	clearTables()

	router := setupRouter()

	// Register the user using the helper function
	w := registerUser("unique_user", "password", router)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Login the user and get the token
	token, w := loginUser("unique_user", "password", router)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, token)

	// Retrieve the registered user from the database
	var user models.User
	config.DB.Where("username = ?", "unique_user").First(&user)

	jsonValue, _ := json.Marshal(map[string]string{
		"Title":   "Test Title",
		"Content": "Test Content",
	})

	req, _ := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the post ID from the response body
	var responseData map[string]interface{}
	json.NewDecoder(w.Body).Decode(&responseData)

	postID := responseData["ID"]

	editJson, _ := json.Marshal(map[string]string{
		"title":   "New Title",
		"content": "New Content",
	})

	editReq, _ := http.NewRequest("PUT", fmt.Sprintf("/api/posts/%.0f", postID), bytes.NewBuffer(editJson))
	editReq.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, editReq)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Add token for authentication in the update request
	editReq.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, editReq)
	assert.Equal(t, http.StatusOK, w.Code)

}
