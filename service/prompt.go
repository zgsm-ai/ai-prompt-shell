package service

import (
	"ai-prompt-shell/dao"
)

var prompts = dao.NewPromptCache()

func Prompts() (map[string]dao.Prompt, error) {
	return prompts.All(), nil
}

func Prompt(prompt_id string) (dao.Prompt, dao.PromptOrigin) {
	return prompts.Get(prompt_id)
}

func PromptIDs() ([]string, error) {
	var result []string
	for k, _ := range prompts.All() {
		result = append(result, k)
	}
	return result, nil
}
