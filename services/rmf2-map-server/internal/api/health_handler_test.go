package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	// Decalre new handler with a temporary directory and test organisation, nil for no VDA5050 poller
	h := NewHandler(t.TempDir(), "test-org", nil)

	// Create a new HTTP request for the /health endpoint - nil means no request body
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Create a fake response writer
	// Instead of sending a real network response, the handler writes into rr, and the test can inspect:
	/*
		rr.Code
		rr.Body
		rr.Header()
	*/
	rr := httptest.NewRecorder()

	// Directly call the handler function with the fake request and response writer
	h.Health(rr, req)

	// This checks the HTTP status code.
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// This creates a variable to hold the decoded JSON response.
	// Since response is expected to be a JSON object with string keys and string values, we use map[string]string.
	/* e.g.
	{
	"status": "ok"
	}
	*/
	var body map[string]string

	// This reads the response body and parses it as JSON. - If the response is not valid JSON, the test fails.
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	// This checks that the "status" field in the JSON response is "ok". If not, the test fails.
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}

}

func TestRoot(t *testing.T) {
	// Decalre new handler with a temporary directory and test organisation, nil for no VDA5050 poller
	h := NewHandler(t.TempDir(), "test-org", nil)

	// Create a new HTTP request for the / endpoint - nil means no request body
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create a fake response writer
	rr := httptest.NewRecorder()

	// Directly call the handler function with the fake request and response writer
	h.Root(rr, req)

	// This checks the HTTP status code.
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// This creates a variable to hold the decoded JSON response.
	var body map[string]string

	// This reads the response body and parses it as JSON. - If the response is not valid JSON, the test fails.
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	// This checks that the "status" field in the JSON response is "ok". If not, the test fails.
	if body["message"] != "rmf2 map backend is running" {
		t.Fatalf("expected message 'rmf2 map backend is running', got %q", body["message"])
	}

}

func TestGetOrganisation(t *testing.T) {
	// Decalre new handler with a temporary directory and test organisation, nil for no VDA5050 poller
	h := NewHandler(t.TempDir(), "test-org", nil)

	// Create a new HTTP request for the /organisation endpoint - nil means no request body
	req := httptest.NewRequest(http.MethodGet, "/organisation", nil)

	// Create a fake response writer
	rr := httptest.NewRecorder()

	// Directly call the handler function with the fake request and response writer
	h.GetOrganisation(rr, req)

	var body map[string]string
	// This reads the response body and parses it as JSON. - If the response is not valid JSON, the test fails.
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	// This checks that the "organisation" field in the JSON response is "test-org". If not, the test fails.
	if body["organisation"] != "test-org" {
		t.Fatalf("expected organisation 'test-org', got %q", body["organisation"])
	}

}
