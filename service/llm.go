package service

import (
	"ai-prompt-shell/dao"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LLMClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// ChatRequest defines chat completion request structure
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []dao.Message `json:"messages"`
}

type ChatResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		LogProbs     struct {
		} `json:"logprobs"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

/**
 * Create new LLM client instance
 * @param baseURL base URL for LLM API endpoint
 * @param apiKey authentication key for API access
 * @return initialized LLM client instance
 */
func NewLLMClient(baseURL, apiKey string) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

/**
 * Execute chat completion using LLM API
 * @param ctx context for request cancellation
 * @param req chat request containing model and messages
 * @return chat completion response from LLM
 * @return error if API call fails
 */
func (c *LLMClient) ChatCompletion(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return ChatResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/v1/chat/completions", c.baseURL),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return ChatResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return ChatResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return ChatResponse{}, fmt.Errorf("LLM API error: %s", resp.Status)
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ChatResponse{}, err
	}

	return result, nil
}
