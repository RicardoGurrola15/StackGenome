package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"stackgenome/pkg/schema/v1"
)

// DefaultEndpoint is the URL of the StackGenome API on Cloudflare Workers
const DefaultEndpoint = "https://stackgenome-api-staging.stackgenome.workers.dev"

// Client is a lightweight HTTP client for the StackGenome remote API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API client with a default timeout of 10 seconds.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultEndpoint
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RecommendationsRequest matches the expected JSON body for POST /v1/recommendations
type RecommendationsRequest struct {
	SchemaVersion string             `json:"schema_version"`
	Fingerprint   FingerprintPayload `json:"fingerprint"`
	Limit         int                `json:"limit"`
}

type FingerprintPayload struct {
	Nodes []schema.NodeDTO `json:"nodes"`
	Edges []schema.EdgeDTO `json:"edges"`
}

// RecommendationsResponse matches the response format from POST /v1/recommendations
type RecommendationsResponse struct {
	RankingVersion  string                     `json:"ranking_version"`
	Recommendations []schema.RecommendationDTO `json:"recommendations"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GetRecommendations sends the sanitized graph to the backend and returns the remote recommendations.
func (c *Client) GetRecommendations(ctx context.Context, dto *schema.ProjectGraphDTO, limit int) ([]schema.RecommendationDTO, error) {
	reqBody := RecommendationsRequest{
		SchemaVersion: "1.0.0",
		Fingerprint: FingerprintPayload{
			Nodes: dto.Nodes,
			Edges: dto.Edges,
		},
		Limit: limit,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/v1/recommendations", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "stackgenome-cli/v1") // Ideally inject Version

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("remote API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		bodyBytes, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error.Message)
		}
		return nil, fmt.Errorf("API error: unexpected status %d", resp.StatusCode)
	}

	var successResp RecommendationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&successResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return successResp.Recommendations, nil
}
