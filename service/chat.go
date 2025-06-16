package service

type ChatModelParameters struct {
	Model     string                 `json:"model"`
	Variables map[string]interface{} `json:"variables"`
}

type ChatModelResponse struct {
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

func ChatWithPrompt(promptId string, model string, variables map[string]interface{}) (ChatModelResponse, error) {
	//TODO:
	return ChatModelResponse{}, nil
}
