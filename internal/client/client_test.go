package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sorolens/sorolens-cli/internal/client"
)

func serve(t *testing.T, handler http.HandlerFunc) *client.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return client.New(srv.URL, srv.Client())
}

func TestGetContract(t *testing.T) {
	want := client.Contract{
		ContractID:      "CTEST",
		Alias:           "my-contract",
		Network:         "mainnet",
		EventCount:      42,
		InvocationCount: 10,
		SuccessRate:     0.95,
	}
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/contracts/CTEST" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})

	got, err := c.GetContract(context.Background(), "CTEST")
	if err != nil {
		t.Fatal(err)
	}
	if got.Alias != want.Alias {
		t.Errorf("alias: got %q, want %q", got.Alias, want.Alias)
	}
	if got.EventCount != want.EventCount {
		t.Errorf("event_count: got %d, want %d", got.EventCount, want.EventCount)
	}
}

func TestGetEvents(t *testing.T) {
	want := client.EventsResponse{
		Events: []client.Event{
			{Time: "2025-01-01T00:00:00Z", Type: "invoke", TxHash: "abc123", Summary: "Called foo()"},
		},
	}
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/contracts/CTEST/events" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})

	got, err := c.GetEvents(context.Background(), "CTEST", "", "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got.Events))
	}
	if got.Events[0].TxHash != "abc123" {
		t.Errorf("tx_hash: got %q, want %q", got.Events[0].TxHash, "abc123")
	}
}

func TestGetEventsQueryParams(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("type") != "invoke" {
			t.Errorf("expected type=invoke, got %q", q.Get("type"))
		}
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %q", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(client.EventsResponse{})
	})

	_, err := c.GetEvents(context.Background(), "CTEST", "invoke", "", 5)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetStorage(t *testing.T) {
	want := client.StorageResponse{
		CurrentLedger: 1000,
		Entries: []client.StorageEntry{
			{KeyHash: "deadbeef", Durability: "persistent", LiveUntilLedger: 1500, LedgersLeft: 500},
		},
	}
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})

	got, err := c.GetStorage(context.Background(), "CTEST")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
}

func TestTrack(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var req client.TrackRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.ContractID != "CTEST" {
			t.Errorf("contract_id: got %q, want %q", req.ContractID, "CTEST")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(client.TrackResponse{ContractID: req.ContractID, Message: "tracked"})
	})

	got, err := c.Track(context.Background(), "CTEST", "")
	if err != nil {
		t.Fatal(err)
	}
	if got.Message != "tracked" {
		t.Errorf("message: got %q, want %q", got.Message, "tracked")
	}
}

func TestAPIError(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"contract not found"}`))
	})

	_, err := c.GetContract(context.Background(), "MISSING")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*client.APIError)
	if !ok {
		t.Fatalf("expected *client.APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("status: got %d, want 404", apiErr.StatusCode)
	}
	if apiErr.Message != "contract not found" {
		t.Errorf("message: got %q", apiErr.Message)
	}
}
