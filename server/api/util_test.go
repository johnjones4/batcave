package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"main/mocks"

	"go.uber.org/mock/gomock"
)

func TestHandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockHookableLogger(ctrl)

	err := errors.New("Sample error")

	logger.EXPECT().Error(err).Times(1)

	// Create a new API instance with a mock logger
	api := &API{APIParams: APIParams{
		Log: logger,
	}}

	// Create a mock HTTP response writer and request
	w := httptest.NewRecorder()

	// Call handleError with a sample error and status
	api.handleError(w, err, http.StatusNotFound)

	// Check if the response has the expected status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Check if the error Message is present in the response body
	expectedResponse := http.StatusText(http.StatusNotFound)
	if !strings.Contains(w.Body.String(), expectedResponse) {
		t.Errorf("Expected response body to contain '%s', got '%s'", expectedResponse, w.Body.String())
	}
}

func TestJsonResponse(t *testing.T) {
	// Create a new API instance
	api := &API{}

	// Create a mock HTTP response writer
	w := httptest.NewRecorder()

	// Create a sample JSON response object
	response := map[string]string{"Message": "Hello, World!"}

	// Call jsonResponse with the sample response
	api.jsonResponse(w, response)

	// Check if the response has the expected content type header
	expectedContentType := "application/json"
	if w.Header().Get("Content-type") != expectedContentType {
		t.Errorf("Expected Content-type header '%s', got '%s'", expectedContentType, w.Header().Get("Content-type"))
	}

	// Check if the response body matches the JSON-encoded response
	expectedResponse, _ := json.Marshal(response)
	if w.Body.String() != string(expectedResponse) {
		t.Errorf("Expected response body '%s', got '%s'", string(expectedResponse), w.Body.String())
	}
}

func TestReadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockHookableLogger(ctrl)

	// Create a new API instance
	api := &API{
		APIParams: APIParams{
			Log: logger,
		},
	}

	// Create a sample JSON request body
	requestBody := `{"name": "John", "age": 30}`

	logger.EXPECT().Debug(requestBody)

	// Create a mock HTTP request with the sample JSON body
	req := httptest.NewRequest("POST", "/", strings.NewReader(requestBody))

	// Create a sample struct to unmarshal the JSON into
	var data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Call readJson to unmarshal the JSON into the data struct
	err := api.readJson(req, &data)

	// Check for any errors during unmarshaling
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the data struct contains the expected values
	if data.Name != "John" || data.Age != 30 {
		t.Errorf("Expected data: {Name: 'John', Age: 30}, got: %+v", data)
	}
}
