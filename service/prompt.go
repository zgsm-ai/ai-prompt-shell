package service

import (
	"ai-prompt-shell/dao"
)

var prompts = dao.NewPromptCache()

/**
 * Get prompt content and origin by ID
 * @param prompt_id ID of the prompt to retrieve
 * @return prompt content
 * @return origin of the prompt (builtin/user-defined)
 */
func Prompt(prompt_id string) (dao.Prompt, dao.PromptOrigin) {
	return prompts.Get(prompt_id)
}

/**
 * Get all available prompt IDs
 * @return slice of prompt IDs
 * @return nil error for future compatibility (always succeeds currently)
 */
func PromptIDs() ([]string, error) {
	var result []string
	for k, _ := range prompts.All() {
		result = append(result, k)
	}
	return result, nil
}
