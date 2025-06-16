package service

import (
	"ai-prompt-shell/dao"
	"ai-prompt-shell/internal/utils"
	"bytes"
	"context"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

type RenderStats struct {
	Total         int64
	Success       int64
	Failures      int64
	TotalDuration time.Duration
	AvgDuration   time.Duration
	CacheHits     int64
	CacheMisses   int64
	RecentErrors  []error
	mu            sync.RWMutex
}

type ToolExecutor func(args ...interface{}) (interface{}, error)

type Renderer struct {
	refreshInterval time.Duration
	templates       map[string]*template.Template
	funcMap         template.FuncMap
	stats           *RenderStats
}

var renderer *Renderer = NewRenderer()

func NewRenderer() *Renderer {
	return &Renderer{
		refreshInterval: 10 * time.Second,
		templates:       make(map[string]*template.Template),
		funcMap:         make(template.FuncMap),
		stats:           &RenderStats{},
	}
}

func constructContextData(variables map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for k, v := range env.All() {
		data[k] = v
	}
	data["variables"] = variables
	return data
}

func onRefreshTools() {
	newFuncs := make(template.FuncMap)
	for k, v := range registry.All() {
		newFuncs[k] = newToolExecutor(&v)
	}
	renderer.funcMap = newFuncs
}

func newToolExecutor(t *dao.Tool) ToolExecutor {
	return func(args ...interface{}) (interface{}, error) {
		return Call(context.Background(), t, args)
	}
}

func onRefreshPrompts() {
	for key, content := range prompts.All() {
		if content.Prompt != "" {
			key = key + ".prompt"
			t, err := template.New(key).Funcs(renderer.funcMap).Parse(content.Prompt)
			if err != nil {
				continue
			}
			renderer.templates[key] = t
		} else if content.Messages != nil {
			for i, _ := range content.Messages {
				msgkey := fmt.Sprintf("%s.messages.%d", key, i)
				t, err := template.New(msgkey).Funcs(renderer.funcMap).Parse(content.Messages[i].Content)
				if err != nil {
					continue
				}
				renderer.templates[msgkey] = t
			}
			continue
		} else {
			logrus.Errorf("prompt %s is invalid", key)
		}
	}
}

func renderTemplate(templateKey string, variables map[string]interface{}) (string, error) {
	t, ok := renderer.templates[templateKey]
	if !ok {
		return "", utils.ErrTemplateNotFound
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, constructContextData(variables))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderMessages(prompt_id string, messages []dao.Message, variables map[string]interface{}) (interface{}, error) {
	var results []dao.Message
	for i, message := range messages {
		content, err := renderTemplate(fmt.Sprintf("%s.messages.%d", prompt_id, i), variables)
		if err != nil {
			return []dao.Message{}, err
		}
		results = append(results, dao.Message{
			Content: content,
			Role:    message.Role,
		})
	}
	return results, nil
}

func RenderPrompt(prompt_id string, variables map[string]interface{}) (string, interface{}, error) {
	prompt, origin := prompts.Get(prompt_id)
	if origin == dao.PromptOrigin_Notexist {
		return "", "", utils.ErrTemplateNotFound
	}
	if prompt.Prompt != "" {
		text, err := renderTemplate(prompt_id+".prompt", variables)
		return "prompt", text, err
	} else if prompt.Messages != nil {
		messages, err := renderMessages(prompt_id, prompt.Messages, variables)
		return "messages", messages, err
	}
	return "", "", utils.ErrTemplateNotFound
}
