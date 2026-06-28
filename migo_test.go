package migo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient()
	cl, ok := c.(*client)
	if !ok {
		t.Fatal("NewClient() should return *client")
	}
	if cl.baseURL != "https://misskey.io/api" {
		t.Errorf("baseURL = %q, want %q", cl.baseURL, "https://misskey.io/api")
	}
	if cl.token != "" {
		t.Errorf("token = %q, want empty", cl.token)
	}
	if cl.httpClient.Timeout != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", cl.httpClient.Timeout)
	}
	if cl.userAgent != "migo (Go Misskey client)" {
		t.Errorf("userAgent = %q", cl.userAgent)
	}
}

func TestNewClient_Options(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	c := NewClient(
		WithBaseURL("https://example.com/api"),
		WithToken("test-token"),
		WithHTTPClient(httpClient),
		WithUserAgent("custom-agent"),
	)
	cl, ok := c.(*client)
	if !ok {
		t.Fatal("NewClient() should return *client")
	}
	if cl.baseURL != "https://example.com/api" {
		t.Errorf("baseURL = %q", cl.baseURL)
	}
	if cl.token != "test-token" {
		t.Errorf("token = %q", cl.token)
	}
	if cl.httpClient != httpClient {
		t.Errorf("httpClient not set")
	}
	if cl.userAgent != "custom-agent" {
		t.Errorf("userAgent = %q", cl.userAgent)
	}
}

func TestClient_Do_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer my-token" {
			t.Errorf("auth header = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("User-Agent") != "test-agent" {
			t.Errorf("user-agent = %q", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Errorf("content-type = %q", r.Header.Get("Content-Type"))
		}

		// Parse body
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["key"] != "value" {
			t.Errorf("body key = %v", body["key"])
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "ok"}`))
	}))
	defer srv.Close()

	client := NewClient(
		WithBaseURL(srv.URL),
		WithToken("my-token"),
		WithUserAgent("test-agent"),
	)

	var resp struct {
		Result string `json:"result"`
	}
	if err := client.Do(context.Background(), "/test", map[string]string{"key": "value"}, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Result != "ok" {
		t.Errorf("result = %q, want %q", resp.Result, "ok")
	}
}

func TestClient_Do_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code": "INVALID_PARAM", "message": "invalid parameter", "id": "abc123"}`))
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	err := client.Do(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d", apiErr.StatusCode)
	}
	if apiErr.Code != "INVALID_PARAM" {
		t.Errorf("Code = %q", apiErr.Code)
	}
	if apiErr.Message != "invalid parameter" {
		t.Errorf("Message = %q", apiErr.Message)
	}
	if apiErr.ID != "abc123" {
		t.Errorf("ID = %q", apiErr.ID)
	}
}

func TestClient_Do_APIError_NonJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	err := client.Do(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d", apiErr.StatusCode)
	}
	if apiErr.Code != "" {
		t.Errorf("Code should be empty for non-JSON error")
	}
}

func TestClient_Do_ContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := client.Do(ctx, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestClient_Do_NilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	if err := client.Do(context.Background(), "/test", nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReqNilRespNotNil(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pong": 1}`))
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	var resp struct {
		Pong int `json:"pong"`
	}
	if err := client.Do(context.Background(), "/ping", nil, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Pong != 1 {
		t.Errorf("Pong = %d, want 1", resp.Pong)
	}
}

func TestNewServices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	svc := NewServices(client)

	if svc.Notes == nil {
		t.Error("Notes is nil")
	}
	if svc.Users == nil {
		t.Error("Users is nil")
	}
	if svc.Drive == nil {
		t.Error("Drive is nil")
	}
	if svc.Auth == nil {
		t.Error("Auth is nil")
	}
	if svc.Timeline == nil {
		t.Error("Timeline is nil")
	}
	if svc.Following == nil {
		t.Error("Following is nil")
	}
	if svc.Notification == nil {
		t.Error("Notification is nil")
	}
	if svc.Admin == nil {
		t.Error("Admin is nil")
	}
	if svc.Chat == nil {
		t.Error("Chat is nil")
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Code:       "SOME_ERROR",
		Message:    "something went wrong",
		ID:         "xyz",
	}
	want := "misskey API error [400] SOME_ERROR: something went wrong (id=xyz)"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
