package service

import (
	"ai-prompt-shell/dao"
)

var env = dao.NewEnvironments()

func Environments() *dao.Environments {
	return env
}
