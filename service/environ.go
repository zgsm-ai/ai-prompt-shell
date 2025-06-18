package service

import (
	"github.com/zgsm-ai/ai-prompt-shell/dao"
)

var environs = dao.NewEnvironments()

/**
 * Get environment variables manager
 * @return pointer to environments instance
 */
func Environments() *dao.Environments {
	return environs
}
