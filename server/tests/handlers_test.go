package tests

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suyashkumar/conduit/server/handlers"
)

func TestMain(m *testing.M) {
	setupTests()
	retCode := m.Run()
	os.Exit(retCode)
}

func setupTests() {
	os.Setenv("DEV", "TRUE") // Set dev environment
	//TODO: Bootstrap dev database?
}

func TestGetUser(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/me", nil)
	w := httptest.NewRecorder()
	sampleEmail := "test@suyash.io"
	samplePrefix := "myPrefix"
	hc := handlers.HomeAutoClaims{Email: sampleEmail, Prefix: samplePrefix}
	context := handlers.Context{}
	handlers.GetUser(w, req, nil, &context, &hc)
	var user handlers.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(t, err)
	assert.Equal(t, user.Email, sampleEmail, "Test emails are equal")
}
