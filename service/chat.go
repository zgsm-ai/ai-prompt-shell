package service

import (
	"ai-prompt-shell/dao"
	"context"
)

/**
 * Execute chat completion with specified prompt template
 * @param promptId ID of the prompt template to use
 * @param model LLM model to use for completion
 * @param args args to substitute into the prompt
 * @return chat response containing completion results
 * @return error if template rendering or LLM call fails
 */
func ChatWithPrompt(promptId string, model string, args map[string]interface{}) (ChatResponse, error) {
	var resp ChatResponse
	// 1. Render template
	kind, data, err := RenderPrompt(promptId, args)
	if err != nil {
		return resp, err
	}

	// 3. Call LLM (Retry 2 times)
	var lastErr error
	var llmReq ChatRequest

	if kind == "prompt" {
		llmReq = ChatRequest{
			Model: model,
			Messages: []dao.Message{
				{
					Role:    "user",
					Content: data.(string),
				},
			},
		}
	} else {
		llmReq = ChatRequest{
			Model:    model,
			Messages: data.([]dao.Message),
		}
	}

	resp, lastErr = llmClient.ChatCompletion(context.Background(), llmReq)

	if lastErr != nil {
		return resp, lastErr
	}
	//TODO:
	return resp, nil
}
