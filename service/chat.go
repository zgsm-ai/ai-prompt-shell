package service

import (
	"ai-prompt-shell/dao"
	"context"
)

type ChatPromptRequest struct {
	Model            string                 `json:"model"`
	Args             map[string]interface{} `json:"args"`
	Temperature      float64                `json:"temperature,omitempty"`
	MaxTokens        int                    `json:"max_tokens,omitempty"`
	TopP             float64                `json:"top_p,omitempty"`
	FrequencyPenalty float64                `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64                `json:"presence_penalty,omitempty"`
	Stop             []string               `json:"stop,omitempty"`
	N                int                    `json:"n,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
	User             string                 `json:"user,omitempty"`
}

/**
 * ChatWithPrompt executes chat completion using specified prompt template
 * @param promptId ID of the prompt template to use
 * @param req chat request parameters containing:
 *      - Model: LLM model to use
 *      - Args: template parameter substitutions
 *      - Temperature: controls randomness of generation
 *      - MaxTokens: maximum tokens to generate
 *      - Other advanced LLM parameters
 * @return ChatResponse containing generated chat response
 * @return error possible errors include:
 *      - prompt template rendering failure
 *      - LLM service call failure
 *      - parameter validation failure
 * Implementation flow:
 * 1. Render prompt template using promptId and Args
 * 2. Construct LLM request parameters
 * 3. Call LLM service to get completion results
 */
func ChatWithPrompt(promptId string, req ChatPromptRequest) (ChatResponse, error) {
	var resp ChatResponse
	// Render template
	kind, data, err := RenderPrompt(promptId, req.Args)
	if err != nil {
		return resp, err
	}

	// Call LLM
	var llmReq ChatRequest = ChatRequest{
		Model:            req.Model,
		Temperature:      req.Temperature,
		MaxTokens:        req.MaxTokens,
		TopP:             req.TopP,
		FrequencyPenalty: req.FrequencyPenalty,
		PresencePenalty:  req.PresencePenalty,
		Stop:             req.Stop,
		N:                req.N,
		Stream:           req.Stream,
		User:             req.User,
	}
	if kind == "prompt" {
		llmReq.Messages = []dao.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: data.(string),
			},
		}
	} else {
		llmReq.Messages = data.([]dao.Message)
	}

	resp, err = llmClient.ChatCompletion(context.Background(), llmReq)
	//TODO:
	return resp, err
}
