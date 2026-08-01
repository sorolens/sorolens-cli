package client

// Contract is the response shape from GET /api/v1/contracts/{id}.
type Contract struct {
	ContractID       string  `json:"contract_id"`
	Alias            string  `json:"alias"`
	Network          string  `json:"network"`
	FirstSeen        string  `json:"first_seen"`
	EventCount       int64   `json:"event_count"`
	InvocationCount  int64   `json:"invocation_count"`
	SuccessRate      float64 `json:"success_rate"`
	StorageEntries   int64   `json:"storage_entries"`
	NearestTTLExpiry int64   `json:"nearest_ttl_expiry"`
}

// Event is one row from GET /api/v1/contracts/{id}/events.
type Event struct {
	Time    string `json:"time"`
	Type    string `json:"type"`
	TxHash  string `json:"tx_hash"`
	Summary string `json:"summary"`
}

// EventsResponse wraps the paginated events list.
type EventsResponse struct {
	Events     []Event `json:"events"`
	NextCursor string  `json:"next_cursor"`
}

// StorageEntry is one row from GET /api/v1/contracts/{id}/storage.
type StorageEntry struct {
	KeyHash         string `json:"key_hash"`
	Durability      string `json:"durability"`
	LiveUntilLedger int64  `json:"live_until_ledger"`
	LedgersLeft     int64  `json:"ledgers_left"`
}

// StorageResponse wraps the storage list.
type StorageResponse struct {
	Entries       []StorageEntry `json:"entries"`
	CurrentLedger int64          `json:"current_ledger"`
}

// TrackRequest is the body for POST /api/v1/contracts.
type TrackRequest struct {
	ContractID string `json:"contract_id"`
	Alias      string `json:"alias,omitempty"`
}

// TrackResponse is the body returned by POST /api/v1/contracts.
type TrackResponse struct {
	ContractID string `json:"contract_id"`
	Alias      string `json:"alias"`
	Message    string `json:"message"`
}

// APIError represents an error response from the API.
type APIError struct {
	StatusCode int
	Message    string `json:"error"`
}

func (e *APIError) Error() string {
	return e.Message
}
